package ioc

import (
	"github.com/LEILEI0628/GinPro/middleware/logger"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

func InitDB(logger loggerx.Logger) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"` // ？此处DSN必须为大写
	}
	var cfg Config
	//var cfg Config = Config{DSN: "root:root@tcp(localhost:3306)/goweb"} // 设置默认值的方法2
	//dsn := viper.GetString("db.mysql.dsn")
	// 使用UnmarshalKey解析"db.mysql"时配置文件必须为树状结构：
	// db:
	//  mysql:
	//   dsn: ""
	// 不能写成db.mysql.dsn: ""的形式（使用配置文件和Remote方式时层级似乎不同，建议直接写成全树状结构）
	err := viper.UnmarshalKey("db.mysql", &cfg) // 需要对cfg进行修改，传指针
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(logger.Debug), glogger.Config{
			// 慢查询阈值：只有执行时间超过这个阈值才打印
			// 一般值为50-100ms，一次磁盘IO一般小于10ms，而SQL查询必然要求命中索引，最好结果是只有一次磁盘IO
			SlowThreshold: time.Millisecond * 50,
			// 是否忽略数据未找到到错误
			IgnoreRecordNotFoundError: true,
			// 日志级别
			LogLevel: glogger.Info,
		}),
	})
	if err != nil {
		// panic相当于整个goroutine结束
		// panic只会出现在初始化的过程中（一旦初始化出错，就没必要启动了）
		panic(err)
	}

	// 建表：
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

type gormLoggerFunc func(msg string, fields ...loggerx.Field)

func (g gormLoggerFunc) Printf(msg string, args ...any) {
	g(msg, loggerx.Field{Key: "args", Value: args})
}
