// facele project main.go
package main

import (
	"flag"

	"log"

	"github.com/hectormao/facele/internal/trns"
	"github.com/hectormao/facele/pkg/cfg"
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

	trns.StartRestAPI(config)
}
