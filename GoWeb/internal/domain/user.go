package domain

import "time"

// User User领域对象（DDD中的聚合根）
// 其他叫法：BO（business object）
type User struct {
	Id         int64
	Email      string
	Phone      string
	Password   string
	NickName   string
	Birthday   string
	AboutMe    string
	CreateTime time.Time
}

type UserProfile struct {
	Email    string
	Phone    string
	NickName string
	Birthday string
	AboutMe  string
}
