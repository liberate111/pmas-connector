package config

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	App AppConfiguration
	DB  DatabaseDriver
	Api ApiConfig
	Log LogConfig
}

type AppConfiguration struct {
	Env       string
	Debug     bool
	InitialDB bool
	DB        string
}

type DatabaseDriver struct {
	Oracle    DatabaseConfiguration
	Mysql     DatabaseConfiguration
	Sqlserver DatabaseConfiguration
}

type DatabaseConfiguration struct {
	Url         string
	Port        int
	ServiceName string
	Username    string
	Password    string
	TableName   string
}

type LogConfig struct {
	Level  string
	Format string
}

type ApiConfig struct {
	Connect RequestApi
	GetData RequestApi
}

type RequestApi struct {
	BasicAuth BasicAuth
	Headers   map[string]string
	Url       string
	Body      string
	Tags      []string
}

type BasicAuth struct {
	Username string
	Password string
}

var Config Configuration

func ReadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %s", err)
	}

	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
