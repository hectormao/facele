package repo

type FacturaCola struct {
	Id  string `json:"id"`
	Xml []byte `json:"xml"`
}

type ColaEnvioRepo interface {
	AgregarEnColaEnvioFacturas(objectId string, facturaXML []byte) error
	GetFacturasAEnviar() (<-chan FacturaCola, error)
}

type ColaNotificacionRepo interface {
	AgregarEnColaNotificacion(objectId string, documentoElectronico []byte) error
	GetFacturasANotificar() (<-chan FacturaCola, error)
}
