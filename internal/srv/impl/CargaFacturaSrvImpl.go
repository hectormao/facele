package impl

import (
	"errors"

	"github.com/hectormao/facele/pkg/repo"
)

type CargaFacturaSrvImpl struct {
	Repo     repo.FacturaRepo
	ColaRepo repo.ColaEnvioRepo
}

var ClienteInvalido = errors.New("Cliente invalido")
var NombreArchivoInvalido = errors.New("Nombre Archivo invalido")
var ContenidoNulo = errors.New("Sin contenido")

func (srv CargaFacturaSrvImpl) Cargar(cliente string, nombreArchivo string, content []byte) (string, error) {
	if cliente == "" {
		return "", ClienteInvalido
	}

	if nombreArchivo == "" {
		return "", NombreArchivoInvalido
	}

	if content == nil {
		return "", ContenidoNulo
	}

	id, err := srv.Repo.AlmacenarFactura(cliente, nombreArchivo, content)
	if err != nil {
		return "", err
	}

	err = srv.ColaRepo.AgregarEnColaEnvioFacturas(id, content)
	if err != nil {
		return "", err
	}

	return id, nil
}
