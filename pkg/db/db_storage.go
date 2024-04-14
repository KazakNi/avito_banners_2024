package db

import (
	"banners/config"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDBConnection() (*sqlx.DB, error) {
	host, port, user, password, dbname, driver := config.Cfg.Database.Host, config.Cfg.Database.Port,
		config.Cfg.Database.User, config.Cfg.Database.Password, config.Cfg.Database.Name, config.Cfg.Database.Driver

	connUrl := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect(driver, connUrl)

	if err != nil {
		panic(fmt.Sprintf("%s, %s", err, connUrl))
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to DB")
	return db, nil
}

var ErrTimeZone = errors.New("failed to load local timezone")

func GetCurrentTime() time.Time {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(ErrTimeZone)
	}
	return time.Now().In(location)
}

func NewDBTestConnection() (*sqlx.DB, error) {
	host, port, user, password, dbname, driver := config.Cfg.Database.Host, config.Cfg.Database.Port,
		config.Cfg.Database.User, config.Cfg.Database.Password, "test", config.Cfg.Database.Driver

	connUrl := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect(driver, connUrl)

	if err != nil {
		panic(fmt.Sprintf("%s, %s", err, connUrl))
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to DB")
	return db, nil
}
