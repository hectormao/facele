package trns

import (
	"encoding/json"

	"github.com/clbanning/mxj"
	facturaSoap "github.com/hectormao/facele/pkg/soap"
)

func AcuseReciboUpdateData(respuesta facturaSoap.AcuseRecibo) (updateData map[string]interface{}) {
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

func XMLToMap(xml []byte) (map[string]interface{}, error) {

	mapValue, err := mxj.NewMapXml(xml)

	if err != nil {
		return nil, err
	}
	return mapValue, nil
}

func StructToMap(data interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonValue, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
