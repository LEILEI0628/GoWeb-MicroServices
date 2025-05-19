package ioc

import (
	"github.com/LEILEI0628/GinPro/middleware/logger"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/conf"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

func InitDB(conf *conf.Bootstrap, logger loggerx.Logger) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"` // ？此处DSN必须为大写
	}
	db, err := gorm.Open(mysql.Open(conf.GetData().GetDatabase().GetSource()), &gorm.Config{
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
