package domain

// User 领域对象, 是DDD中的entity
// BO(business object )
type User struct {
	Id       int64
	Email    string
	Password string
	NickName string
	BirthDay string
	Describe string
}
