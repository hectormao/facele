package main

import (
	"flag"
	"log"

	"github.com/beevik/etree"
	
	"github.com/hectormao/facele/pkg/cfg"
	
	
	srv "github.com/hectormao/facele/internal/srv"
	srvImpl "github.com/hectormao/facele/internal/srv/impl"

	repo "github.com/hectormao/facele/pkg/repo"
	repoImpl "github.com/hectormao/facele/pkg/repo/impl"

	trns "github.com/hectormao/facele/pkg/trns"
	trnsImpl "github.com/hectormao/facele/pkg/trns/impl"
)

func main() {
	gcPath := flag.String("gc", "/home/hectormao/Descargas/lista_valores/Paises-2.1.gc", "GeneriCode path")
	configPath := flag.String("config-file", "../../config/default-config.yaml", "config file path")
	flag.Parse()
	
	config, err := cfg.CargarConfig(*configPath)

	if err != nil {
		log.Fatalf("Error Cargando Config: %v", err)
	}
	
	srv := getServicio(config)
	
	srv.CargarGeneriCodes(gcPath)
	
}

func getServicio(config cfg.FaceleConfigType) srv.CargaGeneriCodeSrv {
	srv := srvImpl.CargaGeneriCodeSrvImpl {
		Repo: repoImpl.GenericodeRepoImpl{
			Config:     config.MongoConfig,
			Translator: trnsImpl.GeneriCodeTrnsImpl{},
		}  
	}
	
	return srv
}
