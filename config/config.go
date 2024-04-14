package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	Database struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host" env:"DB_HOST" env-description:"Database host"`
		Port     string `yaml:"port" env:"DB_PORT" env-description:"Database port"`
		User     string `yaml:"user" env:"DB_USER" env-description:"Database user name"`
		Password string `yaml:"password" env:"DB_PASSWORD" env-description:"Database user password"`
		Name     string `yaml:"dbname" env:"DB_NAME" env-description:"Database name"`
	} `yaml:"database"`
	Server struct {
		Host string `yaml:"host" env:"SRV_HOST,HOST" env-description:"Server host" env-default:"localhost"`
		Port string `yaml:"port" env:"SRV_PORT,PORT" env-description:"Server port" env-default:"8080"`
	} `yaml:"server"`

	Token struct {
		Secret string `yaml:"secret"`
		Salt   string `yaml:"salt"`
	} `yaml:"token"`
}

var Cfg AppConfig

func LoadConfig() {
	err := cleanenv.ReadConfig("./config.yaml", &Cfg)
	if err != nil {
		fmt.Println(err)
		panic("Can't load config data")
	}
}
