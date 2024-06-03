package redisdb

import (
	"github.com/goplateframework/config"

	"github.com/redis/go-redis/v9"
)

func Init(conf *config.Config) (*redis.Client, error) {
	opt, err := redis.ParseURL(conf.RDB.Url)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opt)

	return rdb, nil
}
