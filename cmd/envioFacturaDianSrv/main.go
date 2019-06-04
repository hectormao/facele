// main.go
package main

import (
	"flag"
	"log"

	srv "github.com/hectormao/facele/internal/srv"
	srvImpl "github.com/hectormao/facele/internal/srv/impl"
	"github.com/hectormao/facele/pkg/cfg"
	repo "github.com/hectormao/facele/pkg/repo/impl"
)

func main() {

	configPath := flag.String("config-file", "../../config/default-config.yaml", "config file path")
	flag.Parse()

	log.Printf("Config Path: %v", *configPath)

	config, err := cfg.CargarConfig(*configPath)

	if err != nil {
		log.Fatalf("Error Cargando Config: %v", err)
	}

	log.Printf("Config: %v", config)

	envioFacturaSrv := getServicio(config)
	err = envioFacturaSrv.IniciarConsumidorCola()
	log.Printf("%v", err)
}

func getServicio(config cfg.FaceleConfigType) srv.EnvioFacturaDianSrv {
	facturaRepo := repo.FacturaRepoImpl{Config: config.MongoConfig}
	colaRepo := repo.ColaRepoImpl{Config: config}
	return srvImpl.EnvioFacturaDianSrvImpl{
		facturaRepo,
		colaRepo,
		colaRepo,
		config,
	}
}
