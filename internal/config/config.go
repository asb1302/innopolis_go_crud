package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Config struct {
	AuthServiceHost string
	AuthServiceTLS  bool
}

var config Config

func InitConfig() {
	// TODO добавить энвы включения и отключения джагера
	viper.AutomaticEnv()

	config.AuthServiceHost = viper.GetString("AUTH_SERVICE_HOST")

	authServiceTLS := strings.ToLower(viper.GetString("AUTH_SERVICE_TLS"))
	if authServiceTLS == "true" {
		config.AuthServiceTLS = true
	} else if authServiceTLS == "false" {
		config.AuthServiceTLS = false
	} else {
		log.Fatalf("Invalid value for AUTH_SERVICE_TLS: %s", authServiceTLS)
	}
}

func GetConfig() *Config {
	return &config
}
