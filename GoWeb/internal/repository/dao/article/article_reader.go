package dao

import (
	"context"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao/po"
	"gorm.io/gorm"
)

type ReaderDAO interface {
	Upsert(ctx context.Context, art po.Article) error
	UpsertV2(ctx context.Context, art po.PublishedArticle) error
}

func NewReaderDAO(db *gorm.DB) ReaderDAO {
	panic("implement me")
}
