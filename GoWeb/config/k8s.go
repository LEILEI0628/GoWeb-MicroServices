//go:build k8s

// 使用k8s编译标签

package config

var Config = config{
	DB: DBConfig{
		// 本地连接
		DSN: "root:20010628@tcp(goweb-mysql:13306)/goweb",
	},
	Redis: RedisConfig{
		Addr: "goweb-redis:16379",
	},
}
