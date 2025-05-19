package dao

//
//import (
//	"context"
//	"fmt"
//	"github.com/LEILEI0628/GoWeb/internal/repository/dao/po"
//	"gorm.io/gorm"
//	"time"
//)
//
//func NewArticleDAO(db *gorm.DB) ArticleDAO {
//	return &GORMArticleDAO{
//		db: db,
//	}
//}
//
//func (dao *GORMArticleDAO) Insert(ctx context.Context, article po.Article) (int64, error) {
//	now := time.Now().UnixMilli()
//	article.CreateTime = now
//	article.UpdateTime = now
//	err := dao.db.WithContext(ctx).Create(&article).Error
//	return article.Id, err
//}
//
//func (dao *GORMArticleDAO) UpdateById(ctx context.Context, article po.Article) error {
//	now := time.Now().UnixMilli()
//	article.UpdateTime = now
//	// 依赖GORM忽略零值的特性，默认会用主键进行更新（不推荐，可读性很差）
//	//err := dao.db.WithContext(ctx).Updates(&article).Error
//	res := dao.db.WithContext(ctx).Model(&article).
//		Where("id=? AND author_id=?", article.Id, article.AuthorId). // 防止篡改别人的数据
//		Updates(map[string]any{
//			"title":       article.Title,
//			"content":     article.Content,
//			"update_time": article.UpdateTime,
//		})
//
//	if res.RowsAffected == 0 { // 更新行数
//		return fmt.Errorf("更新失败（可能为创作者非法） id %d, author_id %d",
//			article.Id, article.AuthorId) // 补充日志信息
//	}
//	return res.Error
//}
