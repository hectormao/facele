package trns

import (
	"github.com/hectormao/facele/pkg/ent"
)

type FacturaDianTrns interface {
	FacturaToInvoice(factura ent.FacturaType, resolucion ent.ResolucionFacturacionType) (ent.InvoiceType, error)
}
