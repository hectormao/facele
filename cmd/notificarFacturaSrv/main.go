package main

import (
	"log"

	srv "github.com/hectormao/facele/internal/srv"
	srvImpl "github.com/hectormao/facele/internal/srv/impl"
	repo "github.com/hectormao/facele/pkg/repo/impl"
)

func main() {
	notificarFacturaSrv := getServicio()
	err := notificarFacturaSrv.IniciarConsumidorCola()
	log.Printf("%v", err)
}

func getServicio() srv.NotificarFacturaSrv {
	facturaRepo := repo.FacturaRepoImpl{}
	colaRepo := repo.ColaRepoImpl{}
	return srvImpl.NotificarFacturaSrvImpl{
		facturaRepo,
		colaRepo,
	}
}
