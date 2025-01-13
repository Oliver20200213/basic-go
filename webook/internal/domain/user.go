package domain

import "time"

// User 领域对象, 是DDD中的entity
// BO(business object )
type User struct {
	Id       int64
	Email    string
	Nickname string
	Password string
	Phone    string
	AboutMe  string
	Birthday time.Time
	Ctime    time.Time
}
