//go:build !k8s

// 没有k8s编译标签

package config

var Config = config{
	DB: DBConfig{
		// 本地连接
		DSN: "root:20010628@tcp(localhost:13306)/goweb",
	},
	Redis: RedisConfig{
		Addr: "localhost:16379",
	},
}
