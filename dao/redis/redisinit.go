package redis

import (
	"fmt"
	"bluebell/settings"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", settings.RedisSetting.Host, settings.RedisSetting.Port),
		Password: "",
		DB:       settings.RedisSetting.DB,
		PoolSize: settings.RedisSetting.PoolSize,
	},
	)

	_, err = rdb.Ping().Result()
	return err
}

func Close() {
	_ = rdb.Close()
}
