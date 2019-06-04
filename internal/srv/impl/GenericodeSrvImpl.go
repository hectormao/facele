package impl

import (
	"github.com/hectormao/facele/pkg/ent"
	"github.com/hectormao/facele/pkg/repo"
	"github.com/hectormao/facele/pkg/trns"
)

type GenericodeSrvImpl struct {
	Repo repo.GenericodeRepo
}

func (srv *GenericodeSrvImpl) getCodeName() (map[string]map[string]ent.Genericode, error) {
	codes, err := srv.Repo.GetGenericodesMap()
	if err != nil {
		return nil, err
	}

	return srv.Trans.GenericodeListToMap(codes)

}
