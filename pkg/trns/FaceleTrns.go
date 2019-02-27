package trns

import (
	facturaSoap "github.com/hectormao/facele/pkg/soap"
)

func AcuseReciboApdateData(respuesta facturaSoap.AcuseRecibo) (updateData map[string]interface{}) {
	updateData = make(map[string]interface{})

	updateData["respuestaDian.comments"] = respuesta.Comments
	updateData["respuestaDian.receivedDateTime"] = respuesta.ReceivedDateTime
	updateData["respuestaDian.response"] = respuesta.Response
	updateData["respuestaDian.responseDateTime"] = respuesta.ResponseDateTime
	updateData["respuestaDian.version"] = respuesta.Version

	if respuesta.ReceivedInvoice != nil {
		updateData["respuestaDian.receivedInvoice.comments"] = respuesta.ReceivedInvoice.Comments
		updateData["respuestaDian.receivedInvoice.numeroFactura"] = respuesta.ReceivedInvoice.NumeroFactura
		updateData["respuestaDian.receivedInvoice.response"] = respuesta.ReceivedInvoice.Response
		updateData["respuestaDian.receivedInvoice.UUID"] = respuesta.ReceivedInvoice.UUID
	}

	return
}
