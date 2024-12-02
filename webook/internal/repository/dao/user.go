package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱重复")
	ErrUserNotFund        = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	//存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	//SELECT * FROM users where email=123@qq.com FOR UPDATE
	//查询email=123@qq.com使用FOR UPDATE语句对这些记录进行锁定。
	//这种锁定通常用于确保在多用户环境中，其他事务不能修改或删除这些记录，直到当前事务结束
	//如果在可重复读的隔离级别下，没有改邮箱则会加上间隙锁，如果存在则是加行锁
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			//邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}
func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	//方式二:
	//err := dao.db.WithContext(ctx).First(&u, "email=?",email).Error
	return u, err
	//如果邮箱不存在就会返回gorm.ErrRecordNotFound的错误
}

// User直接对应于数据库表结构
// 有些人叫做 entity,有些人叫做 model 有些人叫做PO(persistent object)
type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	//创建时间，毫秒数
	Ctime int64
	//更新时间，毫秒数
	Utime int64
}
