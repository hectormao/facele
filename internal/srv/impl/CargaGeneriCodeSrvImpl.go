package impl

import (
	"io/ioutil"

	"github.com/hectormao/facele/pkg/repo"
	"github.com/hectormao/facele/pkg/trns"
)

type CargaGeneriCodeSrvImpl struct {
	Repo repo.GenericodeRepo
}

func (srv CargaGeneriCodeSrvImpl) CargarGeneriCodes(filename string) error {

	codes, err := srv.
		Repo.
		GetGenericodes(filename)
	if err != nil {
		return err
	}

	for _, code := range codes {
		srv.
			Repo.
			AlmacenarGenericode(code)
	}

	return nil
}
