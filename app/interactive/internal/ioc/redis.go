package ioc

import (
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/conf"
	"github.com/redis/go-redis/v9"
)

func InitRedis(conf *conf.Bootstrap) redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     conf.GetData().GetRedis().GetAddr(),
		Password: "",
		DB:       0,
	})
}
