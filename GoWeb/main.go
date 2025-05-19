package main

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

// 编译命令：GOOS=linux GOARCH=arm go build -o goweb .
func main() {
	initViper()
	initLogger()
	server := InitWebServer()
	//server.GET("/hello", func(context *gin.Context) {
	//	context.String(http.StatusOK, "hello world")
	//})
	server.Run(":8080")
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	// 不要在日志中打印敏感信息（用户个人信息）
	// msg中应当包含定位信息（确定错误出错的位置）
	// 含糊的日志尽量少使用（如：系统异常）
	zap.S().Debug("启动日志服务（Replace前）") // 不会打印出来
	logger.Debug("使用Logger直接打印")
	zap.ReplaceGlobals(logger)
	zap.S().Debug("启动日志服务")
}

func initViper() {
	//viper.SetConfigName("dev")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath("./config")
	viper.SetConfigFile("./config/dev.yaml")
	// 设置默认值的方法1：viper.SetDefault
	viper.SetDefault("db.mysql.dsn", "root:root@tcp(localhost:3306)/goweb")
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 不太好的点在于只告知了文件变化，没有告知变化的内容
		fmt.Println(in.Name, in.Op)
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperRemote() {
	// etcd写入配置：etcdctl --endpoints=127.0.0.1:12379 put /GoWeb "$(<dev.yaml)"
	// etcd读取配置：etcdctl --endpoints=127.0.0.1:12379 get /GoWeb
	// "/GoWeb"：本项目的配置中心
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:12379", "/GoWeb")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	//err = viper.WatchRemoteConfig()
	//if err != nil {
	//	panic(err)
	//}
	//viper.OnConfigChange(func(in fsnotify.Event) { // 远程配置中心发生变动时不会调用（因为不是fsnotify）
	//
	//})

	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func initViperV1() {
	// 启动时添加参数：--configPath=config/dev.yaml
	// 例：go run . --configPath=config/dev.yaml
	// "config/config.yaml"是默认参数
	cFile := pflag.String("configPath", "config/config.yaml", "配置文件路径")
	pflag.Parse() // 从命令行读取并解析参数
	viper.SetConfigFile(*cFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperReader() { // viper直接读取（io.Reader）
	cfg :=
		`
db:
  mysql:
    dsn: "root:20010628@tcp(localhost:13306)/goweb"
  redis:
    addr: "localhost:16379"
`
	viper.SetConfigType("yaml") // 需要设定配置类型
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}
