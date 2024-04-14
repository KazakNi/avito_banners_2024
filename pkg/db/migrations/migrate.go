package migrations

import (
	"banners/config"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func LoadMigrations() {

	m, err := migrate.New(
		"file://../pkg/db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Cfg.Database.User, config.Cfg.Database.Password,
			config.Cfg.Database.Host, config.Cfg.Database.Port, config.Cfg.Database.Name))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

func LoadTestMigrations() {

	m, err := migrate.New(
		"file://../pkg/db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Cfg.Database.User, config.Cfg.Database.Password,
			config.Cfg.Database.Host, config.Cfg.Database.Port, "test"))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

func TestMigrationsDown() {
	m, err := migrate.New(
		"file://../pkg/db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Cfg.Database.User, config.Cfg.Database.Password,
			config.Cfg.Database.Host, config.Cfg.Database.Port, "test"))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Down(); err != nil {
		log.Fatal(err)
	}

}
