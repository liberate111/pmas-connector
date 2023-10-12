package model

type Configuration struct {
	App  AppConfiguration
	DB   DatabaseDriver
	Log  LogConfig
	Site []SiteConfig
}

type AppConfiguration struct {
	Env       string
	Debug     bool
	InitialDB bool
	Db        string
	Schedule  string
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
}

type LogConfig struct {
	Level  string
	Format string
}

type SiteConfig struct {
	Name      string
	TableName string
	Api       ApiConfig
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
