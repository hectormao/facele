package trns

import (
	"github.com/hectormao/facele/pkg/ent"
)

type FacturaDianTrns interface {
	FacturaToInvoice(factura ent.FacturaType) (ent.InvoiceType, error)
}
