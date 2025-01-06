package domain

import "time"

// User 领域对象, 是DDD中的entity
// BO(business object )
type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	Ctime    time.Time
}
