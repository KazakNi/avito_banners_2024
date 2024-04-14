package main

import (
	"banners/config"
	server "banners/internal"
	"banners/pkg/db/migrations"
)

func main() {

	config.LoadConfig()
	migrations.LoadMigrations()
	server.Run()

}
