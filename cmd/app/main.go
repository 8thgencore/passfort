package main

import (
	"log"
	"os"

	"github.com/8thgencore/passfort/internal/app"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	app.Run(configPath)
}
