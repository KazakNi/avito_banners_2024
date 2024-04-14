package main

import (
	"banners/config"
	server "banners/internal"
)

func main() {
	config.LoadConfig()
	// migrations.LoadMigrations() // убрать и перенести в CLI
	server.Run()

}
