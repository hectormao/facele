// EnvioFacturaDianSrvImpl.go
package impl

import (
	"log"
	"time"

	"bytes"

	"encoding/xml"

	"compress/gzip"

	"github.com/hectormao/facele/pkg/ent"
	"github.com/hectormao/facele/pkg/repo"
	facturaSoap "github.com/hectormao/facele/pkg/soap"
	"github.com/hectormao/facele/pkg/ssl"
	"github.com/hectormao/facele/pkg/trns"
	soap "github.com/hooklift/gowsdl/soap"
	"github.com/satori/go.uuid"
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
		log.Printf("Enviando factura: %v", factura.Id)

		documentoElectronico, err := construirDocumentoElectronico(factura.Xml)

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
			log.Printf("%v", err)
			continue
		}
		if respuesta.Response != 200 {
			log.Printf("Error envio factura DIAN: %v - %v", respuesta.Response, respuesta.Comments)
			srv.ColaNotificacionRepo.AgregarEnColaNotificacion(factura.Id, nil)
		} else {
			log.Printf("Enviando a cola de notificacion: %v", factura.Id)
			srv.ColaNotificacionRepo.AgregarEnColaNotificacion(factura.Id, nil)
		}

		log.Printf("Respuesta Dian: %s", respuesta)

		updateData := trns.AcuseReciboApdateData(*respuesta)

		err = srv.Repo.ActualizarFactura(factura.Id, updateData)
		if err != nil {
			log.Printf("Error actualizando respuesta DIAN: %v", err)
			continue
		}
	}
	return nil
}

func construirDocumentoElectronico(facturaXML []byte) ([]byte, error) {

	invoice, err := ent.NewInvoice()
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

	log.Printf("XML: %s", data)

	var buffer bytes.Buffer

	gz := gzip.NewWriter(&buffer)
	_, err = gz.Write(data)
	if err != nil {
		return nil, err
	}
	if err = gz.Flush(); err != nil {
		return nil, err
	}
	if err = gz.Close(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
