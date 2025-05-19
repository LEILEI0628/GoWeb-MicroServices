package dao

import (
	"context"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/repository/dao/po"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
	GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (po.UserLikeBiz, error)
	DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (po.Interactive, error)
	BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error
	InsertCollectionBiz(ctx context.Context, cb po.UserCollectionBiz) error
	GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (po.UserCollectionBiz, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]po.Interactive, error)
}

type GORMInteractiveDAO struct {
	db *gorm.DB
}

func (dao *GORMInteractiveDAO) GetByIds(ctx context.Context, biz string, ids []int64) ([]po.Interactive, error) {
	var res []po.Interactive
	err := dao.db.WithContext(ctx).Where("biz = ? AND id IN ?", biz, ids).Find(&res).Error
	return res, err
}
func (dao *GORMInteractiveDAO) BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 让调用者保证两者是相等的
		for i := 0; i < len(bizs); i++ {
			err := dao.incrReadCnt(tx, bizs[i], ids[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (dao *GORMInteractiveDAO) incrReadCnt(tx *gorm.DB, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	return tx.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"read_cnt":    gorm.Expr("`read_cnt`+1"),
			"update_time": now,
		}),
	}).Create(&po.Interactive{
		ReadCnt:    1,
		CreateTime: now,
		UpdateTime: now,
		Biz:        biz,
		BizId:      bizId,
	}).Error
}

func (dao *GORMInteractiveDAO) GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (po.UserLikeBiz, error) {
	var res po.UserLikeBiz
	err := dao.db.WithContext(ctx).
		Where("biz=? AND biz_id = ? AND uid = ? AND status = ?",
			biz, bizId, uid, 1).First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (po.UserCollectionBiz, error) {
	var res po.UserCollectionBiz
	err := dao.db.WithContext(ctx).
		Where("biz=? AND biz_id = ? AND uid = ?", biz, bizId, uid).First(&res).Error
	return res, err
}

// InsertCollectionBiz 和 InsertLikeInfo 能不能抽取出来，
// 适当的重复（复制-粘贴）要比强行抽象要更加好一点
// func (dao *GORMInteractiveDAO) common(ctx context.Context,
//
//		biz any, column string, itr Interactive) error {
//		now := time.Now().UnixMilli()
//		return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
//			err := dao.db.WithContext(ctx).Create(&biz).Error
//			if err != nil {
//				return err
//			}
//			return tx.Clauses(clause.OnConflict{
//				DoUpdates: clause.Assignments(map[string]any{
//					column:  gorm.Expr(fmt.Sprintf("`%s`+1", column)),
//					"update_time": now,
//				}),
//			}).Create(&itr).Error
//		})
//	}
//func (dao *GORMInteractiveDAO) InsertCollectionBizV1(ctx context.Context,
//	cb UserCollectionBiz) error {
//	now := time.Now().UnixMilli()
//	return dao.common(ctx, cb, "collect_cnt", Interactive{
//		CollectCnt: 1,
//		CreateTime: now,
//		UpdateTime: now,
//		Biz:        cb.Biz,
//		BizId:      cb.BizId,
//	})
//}

// InsertCollectionBiz 插入收藏记录，并且更新计数
func (dao *GORMInteractiveDAO) InsertCollectionBiz(ctx context.Context,
	cb po.UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	cb.UpdateTime = now
	cb.CreateTime = now
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 插入收藏项目
		err := dao.db.WithContext(ctx).Create(&cb).Error
		if err != nil {
			return err
		}
		// 这边就是更新数量
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"collect_cnt": gorm.Expr("`collect_cnt`+1"),
				"update_time": now,
			}),
		}).Create(&po.Interactive{
			CollectCnt: 1,
			CreateTime: now,
			UpdateTime: now,
			Biz:        cb.Biz,
			BizId:      cb.BizId,
		}).Error
	})
}

func (dao *GORMInteractiveDAO) InsertLikeInfo(ctx context.Context,
	biz string, bizId, uid int64) error {
	// 一把梭
	// 同时记录点赞，以及更新点赞计数
	// 首先你需要一张表来记录，谁点给什么资源点了赞
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先准备插入点赞记录
		// 有没有可能已经点赞过了？
		// 我要不要校验一下，这里必须是没有点赞过
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"update_time": now,
				"statue":      1,
			}),
		}).Create(&po.UserLikeBiz{
			Biz:        biz,
			BizId:      bizId,
			Uid:        uid,
			Status:     1,
			CreateTime: now,
			UpdateTime: now,
		}).Error
		if err != nil {
			return err
		}

		return tx.Clauses(clause.OnConflict{
			// MySQL 不写
			//Columns:
			DoUpdates: clause.Assignments(map[string]any{
				"like_cnt":    gorm.Expr("like_cnt + 1"),
				"update_time": time.Now().UnixMilli(),
			}),
		}).Create(&po.Interactive{
			Biz:        biz,
			BizId:      bizId,
			LikeCnt:    1,
			CreateTime: now,
			UpdateTime: now,
		}).Error
	})
}

func (dao *GORMInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	now := time.Now().UnixMilli()
	// 控制事务超时
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 两个操作
		// 一个是软删除点赞记录
		// 一个是减点赞数量
		err := tx.Model(&po.UserLikeBiz{}).
			Where("biz=? AND biz_id = ? AND uid = ?", biz, bizId, uid).
			Updates(map[string]any{
				"update_time": now,
				"status":      0,
			}).Error
		if err != nil {
			return err
		}
		return tx.Model(&po.Interactive{}).
			// 这边命中了索引，然后没找到，所以不会加锁
			Where("biz=? AND biz_id = ?", biz, bizId).
			Updates(map[string]any{
				"update_time": now,
				"like_cnt":    gorm.Expr("like_cnt-1"),
			}).Error
	})
}

func NewGORMInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{
		db: db,
	}
}

// IncrReadCnt 是一个插入或者更新语义
func (dao *GORMInteractiveDAO) IncrReadCnt(ctx context.Context,
	biz string, bizId int64) error {
	// DAO 要怎么实现？表结构该怎么设计？
	//var itr Interactive
	//err := dao.db.
	//	Where("biz_id =? AND biz = ?", bizId, biz).
	//	First(&itr).Error
	// 两个 goroutine 过来，你查询到 read_cnt 都是 10
	//if err != nil {
	//	return err
	//}
	// 都变成了 11
	//cnt := itr.ReadCnt + 1
	//// 最终变成 11
	//dao.db.Where("biz_id =? AND biz = ?", bizId, biz).Updates(map[string]any{
	//	"read_cnt": cnt,
	//})

	// update a = a + 1
	// 数据库帮你解决并发问题
	// 有一个没考虑到，就是，我可能根本没这一行
	// 事实上这里是一个 upsert 的语义
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		// MySQL 不写
		//Columns:
		DoUpdates: clause.Assignments(map[string]any{
			"read_cnt":    gorm.Expr("read_cnt + 1"),
			"update_time": time.Now().UnixMilli(),
		}),
	}).Create(&po.Interactive{
		Biz:        biz,
		BizId:      bizId,
		ReadCnt:    1,
		CreateTime: now,
		UpdateTime: now,
	}).Error
}

func (dao *GORMInteractiveDAO) Get(ctx context.Context, biz string, bizId int64) (po.Interactive, error) {
	var res po.Interactive
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ?", biz, bizId).
		First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) GetItems() ([]CollectionItem, error) {
	// 不记得构造 JOIN 查询
	var items []CollectionItem
	err := dao.db.Raw("", 1, 2, 3).Find(&items).Error
	return items, err
}

type CollectionItem struct {
	Cid   int64
	Cname string
	BizId int64
	Biz   string
}
