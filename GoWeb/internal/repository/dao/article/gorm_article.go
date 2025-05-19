package dao

import (
	"context"
	"errors"
	"github.com/LEILEI0628/GoWeb/internal/repository/dao/po"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (dao *GORMArticleDAO) GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]po.Article, error) {
	var arts []po.Article
	err := dao.db.WithContext(ctx).Model(&po.Article{}).
		Where("author_id = ?", author).
		Offset(offset).
		Limit(limit).
		Order("update_time DESC").
		Find(&arts).Error
	return arts, err
}

func (dao *GORMArticleDAO) GetPubById(ctx context.Context, id int64) (po.PublishedArticle, error) {
	var pub po.PublishedArticle
	err := dao.db.WithContext(ctx).
		Where("id = ?", id).
		First(&pub).Error
	return pub, err
}

func (dao *GORMArticleDAO) GetById(ctx context.Context, id int64) (po.Article, error) {
	var art po.Article
	err := dao.db.WithContext(ctx).Model(&po.Article{}).
		Where("id = ?", id).
		First(&art).Error
	return art, err
}

func (dao *GORMArticleDAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&po.Article{}).
			Where("id=? AND author_id = ?", id, author).
			Update("status", status)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return ErrPossibleIncorrectAuthor
		}

		res = tx.Model(&po.PublishedArticle{}).
			Where("id=? AND author_id = ?", id, author).Update("status", status)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return ErrPossibleIncorrectAuthor
		}
		return nil
	})
}

func (dao *GORMArticleDAO) Sync(ctx context.Context,
	art po.Article) (int64, error) {
	tx := dao.db.WithContext(ctx).Begin()
	now := time.Now().UnixMilli()
	defer tx.Rollback()
	txDAO := NewGORMArticleDAO(tx)
	var (
		id  = art.Id
		err error
	)
	if id == 0 {
		id, err = txDAO.Insert(ctx, art)
	} else {
		err = txDAO.UpdateById(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	publishArt := po.PublishedArticle(art)
	publishArt.UpdateTime = now
	publishArt.CreateTime = now
	err = tx.Clauses(clause.OnConflict{
		// ID 冲突的时候。实际上，在 MYSQL 里面你写不写都可以
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":       art.Title,
			"content":     art.Content,
			"status":      art.Status,
			"update_time": now,
		}),
	}).Create(&publishArt).Error
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, tx.Error
}

func (dao *GORMArticleDAO) SyncClosure(ctx context.Context,
	art po.Article) (int64, error) {
	var (
		id = art.Id
	)
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		var err error
		now := time.Now().UnixMilli()
		txDAO := NewGORMArticleDAO(tx)
		if id == 0 {
			id, err = txDAO.Insert(ctx, art)
		} else {
			err = txDAO.UpdateById(ctx, art)
		}
		if err != nil {
			return err
		}
		art.Id = id
		publishArt := art
		publishArt.UpdateTime = now
		publishArt.CreateTime = now
		return tx.Clauses(clause.OnConflict{
			// ID 冲突的时候。实际上，在 MYSQL 里面你写不写都可以
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":       art.Title,
				"content":     art.Content,
				"update_time": now,
			}),
		}).Create(&publishArt).Error
	})
	return id, err
}

func (dao *GORMArticleDAO) Insert(ctx context.Context,
	art po.Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.CreateTime = now
	art.UpdateTime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	// 返回自增主键
	return art.Id, err
}

// UpdateById 只更新标题、内容和状态
func (dao *GORMArticleDAO) UpdateById(ctx context.Context,
	art po.Article) error {
	now := time.Now().UnixMilli()
	res := dao.db.Model(&po.Article{}).WithContext(ctx).
		Where("id=? AND author_id = ? ", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":       art.Title,
			"content":     art.Content,
			"status":      art.Status,
			"update_time": now,
		})
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}
