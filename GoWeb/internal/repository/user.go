package repository

import (
	"context"
	"database/sql"
	"github.com/LEILEI0628/GinPro/middleware/cache"
	"github.com/LEILEI0628/GoWeb/internal/domain"
	"github.com/LEILEI0628/GoWeb/internal/repository/cache"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao/po"
	"go.uber.org/zap"
	"time"
)

var ErrUserEmailDuplicated = dao.ErrUserEmailDuplicated
var ErrUserNotFound = dao.ErrUserNotFound
var ErrKeyNotExist = cachex.ErrKeyNotExist

type UserRepository interface {
	FindById(ctx context.Context, id int64) (domain.User, error)
	Create(ctx context.Context, user domain.User) error
	UpdateById(ctx context.Context, id int64, profile domain.UserProfile) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}
type CacheUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CacheUserRepository{dao: dao, cache: cache}
}

func (repo *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 从cache中寻找
	// 获取用户缓存
	uc, err := repo.cache.Get(ctx, id)
	if err == nil {
		// 从cache中找到数据
		zap.L().Debug("Cache Find")
		return uc, err
	}
	//if errors.Is(err, cachex.ErrKeyNotExist) { // 处理缓存未命中：从cache中没找到数据
	// 设置用户缓存：从dao中寻找并写回cache
	up, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	ud := repo.entityToDomain(up)
	err = repo.cache.Set(ctx, ud.Id, ud)
	if err != nil {
		// 缓存Set失败（记录日志做监控即可，为了防止缓存崩溃的可能）
		zap.L().Error("Cache Set Filed", zap.Error(err))
	}
	zap.L().Debug("Cache Set Success")
	return ud, err
	//} // 注释掉此处if语句代表不管缓存发生什么问题都从数据库加载
	// 当缓存发生除ErrKeyNotExist的错误时由两种解决方案：
	// 1.从数据库加载（偶发错误友好），极个别缓存错误可以解决，但当缓存真的崩溃时，要做好兜底保护数据库（大量访问）
	// 2.不加载，默认缓存崩溃，极个别缓存错误也不解决，用户体验较差
	// 面试时选1，极个别缓存错误可以解决，缓存真的崩溃时可以选择数据库限流（基于内存的单机限流）、布尔过滤器
}

func (repo *CacheUserRepository) Create(ctx context.Context, user domain.User) error {
	up := repo.domainToEntity(user)
	// TODO 操作缓存
	return repo.dao.Insert(ctx, up)
}

func (repo *CacheUserRepository) UpdateById(ctx context.Context, id int64, profile domain.UserProfile) error {
	up := po.User{
		NickName: sql.NullString{String: profile.NickName, Valid: profile.NickName != ""},
		Birthday: sql.NullString{String: profile.Birthday, Valid: profile.Birthday != ""},
		AboutMe:  sql.NullString{String: profile.AboutMe, Valid: profile.AboutMe != ""},
	}
	// 更新数据库
	up, err := repo.dao.UpdateById(ctx, id, up)
	if err != nil {
		return err
	}
	ud := repo.entityToDomain(up)
	// 更新cache
	_ = repo.cache.Delete(ctx, id) // 可能会存在键不存在错误（Set也可判断其他错误，故忽略）
	err = repo.cache.Set(ctx, id, ud)
	if err != nil {
		// 缓存Set失败（记录日志做监控即可，为了防止缓存崩溃的可能）
		zap.L().Error("Cache Set Filed", zap.Error(err))
	} else {
		zap.L().Debug("Cache Update Success")
	}
	return nil
}

func (repo *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(user), nil
}

func (repo *CacheUserRepository) domainToEntity(ud domain.User) po.User {
	return po.User{Id: ud.Id,
		Email:      sql.NullString{String: ud.Email, Valid: ud.Email != ""},
		Phone:      sql.NullString{String: ud.Phone, Valid: ud.Phone != ""},
		Password:   ud.Password,
		NickName:   sql.NullString{String: ud.NickName, Valid: ud.NickName != ""},
		Birthday:   sql.NullString{String: ud.Birthday, Valid: ud.Birthday != ""},
		AboutMe:    sql.NullString{String: ud.AboutMe, Valid: ud.AboutMe != ""},
		CreateTime: ud.CreateTime.UnixMilli(),
	}
}

func (repo *CacheUserRepository) entityToDomain(up po.User) domain.User {
	return domain.User{Id: up.Id,
		Email:      up.Email.String,
		Phone:      up.Phone.String,
		Password:   up.Password,
		NickName:   up.NickName.String,
		Birthday:   up.Birthday.String,
		AboutMe:    up.AboutMe.String,
		CreateTime: time.UnixMilli(up.CreateTime)}
}
