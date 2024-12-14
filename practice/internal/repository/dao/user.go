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
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type User struct {
	Id       int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Email    string `gorm:"nique"`
	Password string
	Utime    int64
	Ctime    int64
}

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsError = 1062
		if mysqlErr.Number == uniqueConflictsError {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err

}
