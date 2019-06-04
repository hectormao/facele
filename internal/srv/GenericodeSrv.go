package srv

import (
	"github.com/hectormao/facele/pkg/ent"
)

type GenericodeSrv interface {
	getGenericodes() (map[string]map[string]ent.Genericode, error)
}
