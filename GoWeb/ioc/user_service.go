package ioc

import (
	"github.com/LEILEI0628/GoWeb/internal/repository"
	"github.com/LEILEI0628/GoWeb/internal/service"
	"go.uber.org/zap"
)

func InitUserService(repo repository.UserRepository) service.UserServiceInterface {
	logger, err := zap.NewDevelopment() // 在这里可以为userService单独提供logger（隐私原因）
	if err != nil {
		panic(err)
	}
	return service.NewUserService(repo, logger) // 装饰器模式
}
