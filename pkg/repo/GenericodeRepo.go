package repo

import (
	"github.com/hectormao/facele/pkg/ent"
)

type GenericodeRepo interface {
	GetGenericodes(filename string) []ent.Genericode
	saveGenericode(code map[string]string) (string, error)
	GetGenericodesMap() (map[string]map[string]ent.Genericode, error)
}
