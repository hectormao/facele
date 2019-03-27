package repo

import (
	"github.com/hectormao/facele/pkg/ent"
)

type ColaEnvioRepo interface {
	AgregarEnColaEnvioFacturas(factura ent.FacturaType) error
	GetFacturasAEnviar() (<-chan ent.FacturaType, error)
}

type ColaNotificacionRepo interface {
	AgregarEnColaNotificacion(factura ent.FacturaType) error
	GetFacturasANotificar() (<-chan ent.FacturaType, error)
}
