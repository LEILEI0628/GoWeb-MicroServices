package dao

import (
	"context"
	"errors"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao/po"
)

var ErrPossibleIncorrectAuthor = errors.New("用户在尝试操作非本人数据")

type ArticleDAO interface {
	Insert(ctx context.Context, art po.Article) (int64, error)
	UpdateById(ctx context.Context, art po.Article) error
	GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]po.Article, error)
	GetById(ctx context.Context, id int64) (po.Article, error)
	GetPubById(ctx context.Context, id int64) (po.PublishedArticle, error)
	Sync(ctx context.Context, art po.Article) (int64, error)
	SyncStatus(ctx context.Context, author, id int64, status uint8) error
}
