package main

import "github.com/8thgencore/passfort/internal/app"

const configPath = "./config/config.yaml"

func main() {
	app.Run(configPath)
}
