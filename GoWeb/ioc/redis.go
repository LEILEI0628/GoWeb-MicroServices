package ioc

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		ADDR string `yaml:"addr"`
	}
	var cfg Config
	//addr := viper.GetString("db.redis.addr")
	err := viper.UnmarshalKey("db.redis", &cfg)
	if err != nil {
		panic(fmt.Errorf("初始化Redis失败：%v\n", err))
	}
	return redis.NewClient(&redis.Options{
		Addr:     cfg.ADDR,
		Password: "",
		DB:       0,
	})
}
