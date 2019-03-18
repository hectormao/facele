package repo

type ColaEnvioRepo interface {
	AgregarEnColaEnvioFacturas(factura []byte) error
	GetFacturasAEnviar() (<-chan map[string]interface{}, error)
}

type ColaNotificacionRepo interface {
	AgregarEnColaNotificacion(factura []byte) error
	GetFacturasANotificar() (<-chan map[string]interface{}, error)
}
