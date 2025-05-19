package dao

import (
	"context"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao/po"
	"gorm.io/gorm"
)

type AuthorDAO interface {
	Insert(ctx context.Context, art po.Article) (int64, error)
	UpdateById(ctx context.Context, article po.Article) error
}

func NewAuthorDAO(db *gorm.DB) AuthorDAO {
	panic("implement me")
}
