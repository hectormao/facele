// EnvioFacturaDianSrvImpl.go
package impl

import (
	"log"
	"time"

	"bytes"

	"encoding/base64"
	"encoding/xml"

	"archive/zip"

	"fmt"

	"github.com/hectormao/facele/pkg/ent"
	"github.com/hectormao/facele/pkg/repo"
	facturaSoap "github.com/hectormao/facele/pkg/soap"
	"github.com/hectormao/facele/pkg/ssl"
	"github.com/hectormao/facele/pkg/trns"
	soap "github.com/hooklift/gowsdl/soap"
	"github.com/satori/go.uuid"

	"io/ioutil"

	"encoding/json"
)

type EnvioFacturaDianSrvImpl struct {
	Repo                 repo.FacturaRepo
	ColaEnvioRepo        repo.ColaEnvioRepo
	ColaNotificacionRepo repo.ColaNotificacionRepo
}

func (srv EnvioFacturaDianSrvImpl) IniciarConsumidorCola() error {
	facturasAEnviar, err := srv.ColaEnvioRepo.GetFacturasAEnviar()
	if err != nil {
		return err
	}
	for factura := range facturasAEnviar {
		log.Printf("Enviando factura: %v", factura["_id"])

		idEmpresa := factura["_empresa"].(string)

		documentoElectronico, err := construirDocumentoElectronico(factura)

		if err != nil {
			log.Printf("Error al construir documento electronico: %v", err)
			continue
		}

		client := soap.NewClient(
			"https://facturaelectronica.dian.gov.co/habilitacion/B2BIntegrationEngine/FacturaElectronica",
		)

		u1 := uuid.NewV4()

		securityHeader := facturaSoap.NewWSSSecurityHeader("e34c1826-251c-45a0-bd71-4139fd6db716", "Auditoria9", u1.String(), time.Now(), "")

		log.Printf("Security Header: %v", securityHeader)

		client.AddHeader(securityHeader)

		service := facturaSoap.NewFacturaElectronicaPortName(client)

		var nit facturaSoap.NitType = "123456789"
		var numeroFactura facturaSoap.InvoiceNumberType = "321"

		request := facturaSoap.EnvioFacturaElectronica{
			NIT:           &nit,
			InvoiceNumber: &numeroFactura,
			IssueDate:     time.Now(),
			Document:      documentoElectronico,
		}

		respuesta, err := service.EnvioFacturaElectronica(&request)
		if err != nil {
			log.Printf("Error del web service: %v", err)
			continue
		}

		facturaJson, err := json.Marshal(factura)
		if err != nil {
			log.Printf("Error del al pasar a json: %v", err)
			continue
		}

		if respuesta.Response != 200 {
			log.Printf("Error envio factura DIAN: %v - %v", respuesta.Response, respuesta.Comments)
			srv.ColaNotificacionRepo.AgregarEnColaNotificacion(facturaJson)
		} else {
			log.Printf("Enviando a cola de notificacion: %v", factura)
			srv.ColaNotificacionRepo.AgregarEnColaNotificacion(facturaJson)
		}

		log.Printf("Respuesta Dian: %s", respuesta)

		updateData := trns.AcuseReciboUpdateData(*respuesta)

		err = srv.Repo.ActualizarFactura(factura["_id"].(string), updateData)
		if err != nil {
			log.Printf("Error actualizando respuesta DIAN: %v", err)
			continue
		}
	}
	return nil
}

func construirDocumentoElectronico(factura map[string]interface{}) ([]byte, error) {

	invoice, err := ent.NewInvoice(factura)
	if err != nil {
		return nil, err
	}

	sign, err := ssl.FirmarXML(invoice)
	if err != nil {
		return nil, err
	}

	invoice.AgregarExtension(sign)

	data, err := xml.Marshal(invoice)
	if err != nil {
		return nil, err
	}

	data = []byte(xml.Header + string(data))

	log.Printf("XML: %s", data)

	zipBytes, err := crearZip(data)
	if err != nil {
		return nil, err
	}

	result := toBase64(zipBytes)

	log.Printf("Result %s", result)

	return []byte(result), nil
}

func crearZip(data []byte) ([]byte, error) {
	var buffer bytes.Buffer

	zipWriter := zip.NewWriter(&buffer)

	nombreArchivo := getNombreArchivo()

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

	buffer, _ := base64.StdEncoding.DecodeString(result)

	ioutil.WriteFile("/home/hectormao/base64.txt", []byte(result), 0644)
	ioutil.WriteFile("/home/hectormao/base64.zip", buffer, 0644)

	return result
}

func getNombreArchivo() string {
	return fmt.Sprintf("face_f%010s%010d.xml", "901005166", 81)
}
