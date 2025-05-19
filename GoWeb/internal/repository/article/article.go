package repository

import (
	"context"
	loggerx "github.com/LEILEI0628/GinPro/middleware/logger"
	"github.com/LEILEI0628/GoWeb/internal/domain"
	"github.com/LEILEI0628/GoWeb/internal/repository"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao/article"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao/po"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
)

type ArticleRepository interface {
	// 可以将ArticleRepository拆分为ArticleAuthorRepository和ArticleReaderRepository
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	// Sync 存储并同步数据
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
	//FindById(ctx context.Context, id int64) domain.Article
}

type CachedArticleRepository struct {
	dao      dao.ArticleDAO
	userRepo repository.UserRepository

	// v1 操作两个 DAO
	readerDAO dao.ReaderDAO
	authorDAO dao.AuthorDAO

	// 耦合了 DAO 操作的东西
	// 正常情况下，如果你要在 repository 层面上操作事务
	// 那么就只能利用 db 开始事务之后，创建基于事务的 DAO
	// 或者，直接去掉 DAO 这一层，在 repository 的实现中，直接操作 db
	db *gorm.DB
	l  loggerx.Logger
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}
func (c *CachedArticleRepository) SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error {
	return c.dao.SyncStatus(ctx, id, author, uint8(status))
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Sync(ctx, c.toEntity(art))
}

//func (c *CachedArticleRepository) SyncV2_1(ctx context.Context, art domain.Article) (int64, error) {
//	// 谁在控制事务，是 repository，还是DAO在控制事务？
//	c.dao.Transaction(ctx, func(txDAO dao.ArticleDAO) error {
//
//	})
//}

// SyncV2 尝试在 repository 层面上解决事务问题
// 确保保存到制作库和线上库同时成功，或者同时失败
func (c *CachedArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	// 开启了一个事务
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer tx.Rollback()
	// 利用 tx 来构建 DAO
	author := dao.NewAuthorDAO(tx)
	reader := dao.NewReaderDAO(tx)

	var (
		id  = art.Id
		err error
	)
	artn := c.toEntity(art)
	// 应该先保存到制作库，再保存到线上库
	if id > 0 {
		err = author.UpdateById(ctx, artn)
	} else {
		id, err = author.Insert(ctx, artn)
	}
	if err != nil {
		// 执行有问题，要回滚
		//tx.Rollback()
		return id, err
	}
	// 操作线上库了，保存数据，同步过来
	// 考虑到，此时线上库可能有，可能没有，你要有一个 UPSERT 的写法
	// INSERT or UPDATE
	// 如果数据库有，那么就更新，不然就插入
	err = reader.UpsertV2(ctx, po.PublishedArticle(artn))
	// 执行成功，直接提交
	tx.Commit()
	return id, err

}

func (c *CachedArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	artn := c.toEntity(art)
	// 应该先保存到制作库，再保存到线上库
	if id > 0 {
		err = c.authorDAO.UpdateById(ctx, artn)
	} else {
		id, err = c.authorDAO.Insert(ctx, artn)
	}
	if err != nil {
		return id, err
	}
	// 操作线上库了，保存数据，同步过来
	// 考虑到，此时线上库可能有，可能没有，你要有一个 UPSERT 的写法
	// INSERT or UPDATE
	// 如果数据库有，那么就更新，不然就插入
	err = c.readerDAO.Upsert(ctx, artn)
	return id, err
}

func (c *CachedArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Insert(ctx, po.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}

func (c *CachedArticleRepository) Update(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateById(ctx, po.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}

func (repo *CachedArticleRepository) GetPublishedById(
	ctx context.Context, id int64) (domain.Article, error) {
	// 读取线上库数据，如果你的 Content 被你放过去了 OSS 上，你就要让前端去读 Content 字段
	art, err := repo.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	// 你在这边要组装 user 了，适合单体应用
	usr, err := repo.userRepo.FindById(ctx, art.AuthorId)
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id:   usr.Id,
			Name: usr.NickName,
		},
		CreateTime: time.UnixMilli(art.CreateTime),
		UpdateTime: time.UnixMilli(art.UpdateTime),
	}
	return res, nil
}

func (c *CachedArticleRepository) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	data, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return c.toDomain(data), nil
}

func (c *CachedArticleRepository) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	// TODO 集成缓存方案
	// 你只缓存这一页
	res, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	data := slice.Map[po.Article, domain.Article](res, func(idx int, src po.Article) domain.Article {
		return c.toDomain(src)
	})
	// TODO 回写缓存
	return data, nil
}
func (c *CachedArticleRepository) toEntity(art domain.Article) po.Article {
	return po.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   uint8(art.Status),
	}
}

func (repo *CachedArticleRepository) toDomain(art po.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		CreateTime: time.UnixMilli(art.CreateTime),
		UpdateTime: time.UnixMilli(art.UpdateTime),
	}
}
