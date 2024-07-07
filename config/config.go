package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

func Load() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath("./config/priv")
	v.SetConfigName("configcloud")
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func Parse(v *viper.Viper) (*Config, error) {
	c := new(Config)

	if err := v.Unmarshal(&c); err != nil {
		log.Printf("unable to decode %v", err)
		return nil, err
	}

	return c, nil
}
