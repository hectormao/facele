package trns

import (
	"github.com/hectormao/facele/pkg/ent"
)

type GeneriCodeTrans interface {
	XMLToGenericode(xmlfile []byte) ([]ent.Genericode, error)
	GenericodeListToMap(genericodes []ent.Genericode) map[string]map[string]ent.Genericode
}
