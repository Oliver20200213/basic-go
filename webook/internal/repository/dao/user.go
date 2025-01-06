package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱重复")
	ErrUserNotFund   = gorm.ErrRecordNotFound
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
		// 邮箱冲突 or 手机号码冲突
		if mysqlErr.Number == uniqueConflictsErrNo {
			//邮箱冲突
			return ErrUserDuplicate
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

func (dao *UserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).First(&u).Error
	return u, err
}

// User直接对应于数据库表结构
// 有些人叫做 entity,有些人叫做 model 有些人叫做PO(persistent object)
type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	// Phone    string `gorm:"unique"`
	// 如果都是用email注册登录，没有手机，那么Phone为空字符串，就会出现unique索引的冲突,
	// 反之都是用phone登录，email为空字符串，也会出现冲突
	// 唯一索引允许有多个空值
	// 但是不能有多个""空字符串
	Phone sql.NullString `gorm:"unique"`

	// 早期有使用该方法，最大的问题是要解引用
	//  需要判空，如果对null进行解引用（*phone）就会触发panic
	// Phone *string

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("`id` = ?", id).First(&u).Error
	return u, err
}
