package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	DB     PostgresConfig
	RDB    RedisConfig
	Logger Logger
}

type ServerConfig struct {
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

type Logger struct {
	Development bool
	Encoding    string
	Filepath    string
	Level       string
}

type PostgresConfig struct {
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

type RedisConfig struct {
	Url string
}

func LoadConfig() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	if err := v.Unmarshal(&c); err != nil {
		log.Printf("unable to decode %v", err)
		return nil, err
	}

	return &c, nil
}
