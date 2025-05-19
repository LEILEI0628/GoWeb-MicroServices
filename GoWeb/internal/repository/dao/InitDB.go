package dao

import (
	"github.com/LEILEI0628/GoWeb/internal/repository/dao/po"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&po.User{}, &po.Article{}, &po.PublishedArticle{})
}
