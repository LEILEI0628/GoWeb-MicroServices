package po

import "database/sql"

// User dao.User直接对应数据库表
// 其他叫法：entity，model，PO（persistent object）
type User struct {
	Id         int64          `gorm:"primaryKey,auto_Increment"` // 自增主键
	Email      sql.NullString `gorm:"unique"`                    // 唯一索引允许有多个空值（null），但是不能有多个空字符串（""）
	Phone      sql.NullString `gorm:"unique"`                    // *string也可以，但是要解引用，判空
	Password   string
	NickName   sql.NullString
	Birthday   sql.NullString
	AboutMe    sql.NullString
	CreateTime int64 // 创建时间：毫秒数
	UpdateTime int64 // 修改时间：毫秒数

}
