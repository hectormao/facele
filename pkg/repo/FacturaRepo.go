// FacturaRepo.go
package repo

type FacturaRepo interface {
	AlmacenarFactura(cliente string, nombreArchivo string, contenido []byte) (string, error)
	ActualizarFactura(id string, data map[string]interface{}) error
}
