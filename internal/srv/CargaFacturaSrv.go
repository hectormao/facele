package srv

type CargaFacturaSrv interface {
	Cargar(empresa string, nombreArchivo string, content []byte) (string, error)
}
