package srv

type CargaFacturaSrv interface {
	Cargar(nombreArchivo string, content []byte) (string, error)
}
