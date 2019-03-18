package impl

import (
	"errors"

	"github.com/hectormao/facele/pkg/repo"
	"github.com/hectormao/facele/pkg/trns"

	"encoding/json"

	"log"
)

type CargaFacturaSrvImpl struct {
	Repo     repo.FacturaRepo
	ColaRepo repo.ColaEnvioRepo
}

var ClienteInvalido = errors.New("Cliente invalido")
var NombreArchivoInvalido = errors.New("Nombre Archivo invalido")
var ContenidoNulo = errors.New("Sin contenido")
var EmpresaInvalida = errors.New("Numero documento empresa invalido")

func (srv CargaFacturaSrvImpl) Cargar(documentoEmpresa string, nombreArchivo string, content []byte) (string, error) {
	if documentoEmpresa == "" {
		return "", ClienteInvalido
	}

	if nombreArchivo == "" {
		return "", NombreArchivoInvalido
	}

	if content == nil {
		return "", ContenidoNulo
	}

	empresa, err := srv.Repo.GetEmpresa(documentoEmpresa)

	if err != nil {
		return "", err
	}

	if empresa == nil {
		return "", EmpresaInvalida
	}

	log.Printf("Se obtiene la empresa: %v", empresa)

	factura, err := trns.XMLToMap(content)
	if err != nil {
		return "", err
	}

	factura["_empresa_id"] = empresa.ObjectId.Hex()

	id, err := srv.Repo.AlmacenarFactura(factura)
	if err != nil {
		return "", err
	}

	factura["_id"] = id
	empresaMap, err := trns.StructToMap(empresa)
	if err != nil {
		return "", err
	}
	factura["_empresa"] = empresaMap

	facturaJson, err := json.Marshal(factura)
	if err != nil {
		return "", err
	}

	err = srv.ColaRepo.AgregarEnColaEnvioFacturas(facturaJson)
	if err != nil {
		return "", err
	}

	return id, nil
}
