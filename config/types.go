package config

import "time"

type Config struct {
	Server        serverConfig
	DB            postgresConfig
	Cache         redisConfig
	Logger        loggerConfig
	GRPCWorker    grpcWorkerConfig
	GoogleStorage googleStorageConfig
}

type serverConfig struct {
	AllowedOrigins        []string
	AppVersion            string
	CookieName            string
	CSRF                  bool
	CtxDefaultTimeout     time.Duration
	Debug                 bool
	Host                  string
	JWTRefreshTokenSecret string
	JWTAccessTokenSecret  string
	Mode                  string
	Port                  string
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
}

type loggerConfig struct {
	Development bool
	Encoding    string
	Filepath    string
	Level       string
}

type postgresConfig struct {
	ConnMaxIddleTime int8
	ConnMaxLifetime  int8
	Dbname           string
	Driver           string
	Host             string
	MaxIddleConns    int
	MaxOpenConns     int
	Password         string
	Port             string
	SSLMode          string
	User             string
	EndpointID       string
}

type grpcWorkerConfig struct {
	Port string
}

type redisConfig struct {
	Url string
}

type googleStorageConfig struct {
	Path       string
	BucketName string
}
