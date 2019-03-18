// facele project main.go
package main

import (
	"flag"

	"log"

	"github.com/hectormao/facele/internal/trns"
	"github.com/hectormao/facele/pkg/cfg"

	"os"
)

func main() {
	configPath := flag.String("config-file", "../../config/default-config.yaml", "config file path")
	flag.Parse()

	log.Printf("Config Path: %v", *configPath)

	dir, _ := os.Getwd()

	log.Printf("Directorio Actual: %v", dir)

	config, err := cfg.CargarConfig(*configPath)

	if err != nil {
		log.Fatalf("Error Cargando Config: %v", err)
	}

	log.Printf("Config: %v", config)

	trns.StartRestAPI(config)
}
