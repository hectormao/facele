package srv

type CargaFacturaSrv interface {
	Cargar(cliente string, nombreArchivo string, content []byte) (string, error)
}
