// FacturaRepo.go
package repo

import (
	"github.com/hectormao/facele/pkg/ent"
)

type FacturaRepo interface {
	AlmacenarFactura(factura ent.FacturaType) (string, error)
	ActualizarFactura(id string, data map[string]interface{}) error
	GetEmpresa(numeroDocumento string) (*ent.EmpresaType, error)
	GetEmpresaPorId(id string) (*ent.EmpresaType, error)
}
