package service

import (
	"context"
	"errors"
	"github.com/LEILEI0628/GoWeb/internal/domain"
	"github.com/LEILEI0628/GoWeb/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserEmailDuplicated = repository.ErrUserEmailDuplicated
var ErrInvalidEmailOrPassword = errors.New("invalid email or password")
var ErrUserNotFound = repository.ErrUserNotFound

type UserServiceInterface interface {
	SignUp(ctx context.Context, user domain.User) error
	SignIn(ctx context.Context, user domain.User) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.UserProfile, error)
	Edit(ctx context.Context, id int64, user domain.UserProfile) error
}
type UserService struct {
	userRepo repository.UserRepository
	logger   *zap.Logger // 预留了注入空间
}

func NewUserService(userRepo repository.UserRepository, logger *zap.Logger) UserServiceInterface {
	return &UserService{userRepo: userRepo, logger: logger}
}

func (service *UserService) SignUp(ctx context.Context, user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return service.userRepo.Create(ctx, user)
}

func (service *UserService) SignIn(ctx context.Context, user domain.User) (domain.User, error) {
	// 通过邮箱查找用户
	userFind, err := service.userRepo.FindByEmail(ctx, user.Email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidEmailOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(userFind.Password), []byte(user.Password))
	if err != nil {
		return domain.User{}, ErrInvalidEmailOrPassword
	}
	return userFind, nil
}

func (service *UserService) Edit(ctx context.Context, id int64, profile domain.UserProfile) error {
	return service.userRepo.UpdateById(ctx, id, profile)
}

func (service *UserService) Profile(ctx context.Context, id int64) (domain.UserProfile, error) {
	user, err := service.userRepo.FindById(ctx, id)
	if err != nil {
		return domain.UserProfile{}, err
	}
	return domain.UserProfile{Email: user.Email, Phone: user.Phone, NickName: user.NickName,
		Birthday: user.Birthday, AboutMe: user.AboutMe}, err
}
