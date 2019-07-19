// EnvioFacturaDianSrvImpl.go
package impl

import (
	"log"
	"time"

	"bytes"

	"strconv"

	"encoding/base64"
	"encoding/xml"

	"archive/zip"
	"errors"
	"fmt"

	"github.com/hectormao/facele/pkg/cfg"
	"github.com/hectormao/facele/pkg/ent"
	"github.com/hectormao/facele/pkg/repo"
	facturaSoap "github.com/hectormao/facele/pkg/soap"
	"github.com/hectormao/facele/pkg/trns"
	soap "github.com/hooklift/gowsdl/soap"
	"github.com/satori/go.uuid"

	"io/ioutil"
)

type EnvioFacturaDianSrvImpl struct {
	Repo                 repo.FacturaRepo
	ColaEnvioRepo        repo.ColaEnvioRepo
	ColaNotificacionRepo repo.ColaNotificacionRepo
	Config               cfg.FaceleConfigType
	FacturaTrns          trns.FacturaDianTrns
}

var NoResolucionesError error = errors.New("Resolucion Activa no existente")

func (srv EnvioFacturaDianSrvImpl) IniciarConsumidorCola() error {
	facturasAEnviar, err := srv.ColaEnvioRepo.GetFacturasAEnviar()
	if err != nil {
		return err
	}

	for factura := range facturasAEnviar {
		log.Printf("Enviando factura: %v", factura.ObjectId)

		resoluciones := factura.Empresa.Resoluciones
		var resolucion *ent.ResolucionFacturacionType
		for _, res := range resoluciones {
			if res.SiActivo {
				resolucion = &res
				break
			}
		}
		if resolucion == nil {
			log.Printf("%v", NoResolucionesError)
			continue
		}

		idEmpresa := factura.EmpresaID
		log.Printf("idEmpresa: %v", idEmpresa)
		factura.Resolucion = *resolucion
		documentoElectronico, cufe, err := srv.construirDocumentoElectronico(factura)

		if err != nil {
			log.Printf("Error al construir documento electronico: %v", err)
			continue
		}

		client := soap.NewClient(
			srv.Config.WebServiceDian.Url,
		)

		u1 := uuid.NewV4()

		securityHeader := facturaSoap.NewWSSSecurityHeader(
			factura.Empresa.SoftwareFacturacion.Id,
			factura.Empresa.SoftwareFacturacion.Password,
			u1.String(),
			time.Now().UTC(),
			"",
		)

		log.Printf("Security Header: %v", *securityHeader)

		client.AddHeader(securityHeader)

		service := facturaSoap.NewFacturaElectronicaPortName(client)

		var nit facturaSoap.NitType = facturaSoap.NitType(factura.CabezaFactura.Nit)
		var numeroFactura facturaSoap.InvoiceNumberType = facturaSoap.InvoiceNumberType(
			resolucion.Prefijo + strconv.Itoa(factura.CabezaFactura.Consecutivo),
		)

		request := facturaSoap.EnvioFacturaElectronica{
			NIT:           &nit,
			InvoiceNumber: &numeroFactura,
			IssueDate:     time.Now(),
			Document:      documentoElectronico,
		}

		xmlRequest, err := xml.Marshal(request)

		log.Printf("DIAN Request: %s", xmlRequest)

		respuesta, err := service.EnvioFacturaElectronica(&request)
		if err != nil {
			log.Printf("Error del web service: %v", err)
			continue
		}

		factura.Cufe = cufe

		if respuesta.Response != 200 {
			log.Printf("Error envio factura DIAN: %v - %v", respuesta.Response, respuesta.Comments)
			srv.ColaNotificacionRepo.AgregarEnColaNotificacion(factura)
		} else {
			log.Printf("Enviando a cola de notificacion: %v", factura)
			srv.ColaNotificacionRepo.AgregarEnColaNotificacion(factura)
		}

		log.Printf("Respuesta Dian: %s", respuesta)

		updateData := trns.AcuseReciboUpdateData(*respuesta)

		err = srv.Repo.ActualizarFactura(factura.ObjectId, updateData)
		if err != nil {
			log.Printf("Error actualizando respuesta DIAN: %v", err)
			continue
		}
	}
	return nil
}

func (srv EnvioFacturaDianSrvImpl) construirDocumentoElectronico(factura ent.FacturaType) ([]byte, string, error) {

	invoice, err := srv.FacturaTrns.FacturaToInvoice(factura)

	data, err := xml.Marshal(invoice)
	if err != nil {
		return nil, "", err
	}

	data = []byte(xml.Header + string(data))

	log.Printf("XML: %s", data)

	zipBytes, err := crearZip(
		data,
		factura.CabezaFactura.Nit,
		strconv.Itoa(factura.CabezaFactura.Consecutivo),
	)
	if err != nil {
		return nil, "", err
	}

	ioutil.WriteFile("/tmp/archivo_enviar.zip", zipBytes, 0644)

	result := toBase64(zipBytes)

	log.Printf("Result %s", result)

	return []byte(result), invoice.UUID.Data, nil
}

func crearZip(data []byte, nitResponsable string, idFactura string) ([]byte, error) {
	var buffer bytes.Buffer

	zipWriter := zip.NewWriter(&buffer)

	nombreArchivo := getNombreArchivo(nitResponsable, idFactura)

	zipFile, err := zipWriter.Create(nombreArchivo)
	if err != nil {
		return nil, err
	}

	_, err = zipFile.Write(data)
	if err != nil {
		return nil, err
	}

	if err = zipWriter.Close(); err != nil {
		return nil, err
	}

	ioutil.WriteFile("/home/hectormao/antes.zip", buffer.Bytes(), 0644)

	return buffer.Bytes(), nil
}

func toBase64(data []byte) string {

	log.Printf("toBase64 data: %d", len(data))

	result := base64.StdEncoding.EncodeToString(data)

	return result
}

func getNombreArchivo(nitResponsable string, idFactura string) string {
	return fmt.Sprintf("face_f%010s%010s.xml", nitResponsable, idFactura)
}
