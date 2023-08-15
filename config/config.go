package config

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	App     AppConfiguration
	MongoDb DatabaseConfiguration
}
type AppConfiguration struct {
	Env   string
	Port  int
	Debug bool
}
type DatabaseConfiguration struct {
	Connection string
}

var Config Configuration

func ReadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
