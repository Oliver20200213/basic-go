package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
	"database/sql"
	"log"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFund
)

// UserRepository 是核心，它有不同的实现
// 但是Factory本身如果只是初始化一下，那么它不是你的核心
type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	UpdateNonZeroFields(ctx context.Context, u domain.User) error
}

type CacheUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))

	//create的时候可以用下缓存，来提高性能，一般注册之后都会登录
	//单独在login的时候没有必要使用缓存
}

func (r *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	//先从cache里面找
	//再从dao里面找
	//找到了回写cache

	u, err := r.cache.Get(ctx, id)
	// 缓存里面有数据
	// 缓存里面没有数据
	// 缓存出错了，你也不知道有没有数据
	if err == nil {
		// 必然是有数据
		return u, nil
	}
	//没这个数据
	//if err == cache.ErrKeyNotExist {
	//	//去数据库里面加载
	//}
	// 除了上面的情况，其他的情况怎么办， err = io.EOF
	// 要不要去数据库加载
	// 加载对偶发性的错误很友好，万一redis崩了加载会导致数据库也崩了

	// 选加载---做好兜底，万一Redis真的崩了，要保护好你的数据库（面试选加载）
	// 数据库进行限流（内存限流，redis已经崩了）

	// 选不加载---用户体验差一点

	// cache中没有找到数据，并且如果是Redis崩了的情况下去数据库中加载
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = r.entityToDomain(ue)
	// 将数据写入到cache中，可以开个goroutine
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 这里该怎么办？
			// 写入日志做监控即可
			log.Println(err)
		}
	}()
	return u, err
}

// 注意：使用缓存的两个核心问题：第一个是一致性问题， 第二个是我的缓存崩了

func (r *CacheUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			// 我确实有手机号
			Valid: u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Nickname: u.Nickname,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday.UnixMilli(), //将time.Time转为int64
		Ctime:    u.Ctime.UnixMilli(),
	}
}

func (r *CacheUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String, //u.Email.String表示存储的值 u.Email.Valid表示有没有数据
		Password: u.Password,
		Phone:    u.Phone.String,
		Nickname: u.Nickname,
		AboutMe:  u.AboutMe,
		Birthday: time.UnixMilli(u.Birthday), //将int64转为time.Time
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

func (r *CacheUserRepository) UpdateNonZeroFields(ctx context.Context, u domain.User) error {
	return r.dao.UpdateById(ctx, r.domainToEntity(u))
}
