package config

import "github.com/ilyakaznacheev/cleanenv"

type AppConfig struct {
}

var Config AppConfig

func LoadConfig() {
	err := cleanenv.ReadConfig("config.yaml", &Config)
	if err != nil {
		panic("Can't load config data")
	}
}
