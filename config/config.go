package config

import (
	"fmt"
	"os"
	"reflect"
)

type AppConfig struct {
	Host string `env:"APP_HOST"`
	Port string `env:"APP_PORT"`
}

type DbConfig struct {
	Host      string `env:"POSTGRES_HOST"`
	Port      string `env:"POSTGRES_PORT"`
	User      string `env:"POSTGRES_USER"`
	Password  string `env:"POSTGRES_PASSWORD"`
	DbName    string `env:"POSTGRES_DB"`
	DbSslMode string `env:"DB_SSLMODE"`
}

type ServiceConfig struct {
	App AppConfig
	Db  DbConfig
}

func envsToStruct(s interface{}) error {
	structValue := reflect.ValueOf(s).Elem()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Type().Field(i)
		envTag := field.Tag.Get("env")

		if envTag == "" {
			return fmt.Errorf("struct has no env tag")
		}

		envVal := os.Getenv(envTag)
		if envVal == "" {
			return fmt.Errorf("env var \"%v\" is empty", envTag)
		}
		structValue.Field(i).SetString(envVal)
	}

	return nil
}

func GetConfig() (ServiceConfig, error) {
	var cfg ServiceConfig

	var appCfg AppConfig
	err := envsToStruct(&appCfg)
	if err != nil {
		return cfg, err
	}

	var dbCfg DbConfig
	err = envsToStruct(&dbCfg)
	if err != nil {
		return cfg, err
	}

	cfg.App = appCfg
	cfg.Db = dbCfg

	return cfg, err
}
