package dao

import (
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/repository/dao/po"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&po.Interactive{}, &po.Collection{}, &po.UserCollectionBiz{}, &po.UserLikeBiz{})
}
