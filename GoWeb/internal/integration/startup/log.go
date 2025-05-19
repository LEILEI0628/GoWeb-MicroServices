package startup

import (
	loggerx "github.com/LEILEI0628/GinPro/middleware/logger"
)

func InitLog() loggerx.Logger {
	return &loggerx.NoneLogger{}
}
