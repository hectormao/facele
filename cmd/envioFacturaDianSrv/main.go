// main.go
package main

import (
	"log"

	srv "github.com/hectormao/facele/internal/srv"
	srvImpl "github.com/hectormao/facele/internal/srv/impl"
	repo "github.com/hectormao/facele/pkg/repo/impl"
)

func main() {

	envioFacturaSrv := getServicio()
	err := envioFacturaSrv.IniciarConsumidorCola()
	log.Printf("%v", err)
}

func getServicio() srv.EnvioFacturaDianSrv {
	facturaRepo := repo.FacturaRepoImpl{}
	colaRepo := repo.ColaRepoImpl{}
	return srvImpl.EnvioFacturaDianSrvImpl{
		facturaRepo,
		colaRepo,
		colaRepo,
	}
}
