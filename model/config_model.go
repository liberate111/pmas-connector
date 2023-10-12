package model

type Configuration struct {
	App       AppConfiguration
	DB        DatabaseDriver
	Log       LogConfig
	Api       ApiConfig
	TableName string
}

type AppConfiguration struct {
	Env       string
	Debug     bool
	InitialDB bool
	Db        string
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
	// TableName   string
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
