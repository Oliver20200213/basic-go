package repository

import (
	"basic-go/practice/internal/domain"
	"basic-go/practice/internal/repository/cache"
	"basic-go/practice/internal/repository/dao"
	"context"
	"database/sql"
	"log"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &CacheUserRepository{dao: dao, cache: cache}
}

func (repo *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.domainToEntity(u))

}

func (repo *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(u), nil

}

func (repo *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 先从cache中查找
	// 再从dao中查找
	// 找到后回写cache

	u, err := repo.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}

	// 如果cache中没有则从数据库中加载（做好兜底，万一Redis真的崩了，要保护好你的数据库）
	ue, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = repo.entityToDomain(ue)
	//将输入会写到cache
	go func() {
		err = repo.cache.Set(ctx, id, u)
		if err != nil {
			// 写入日志进行监控
			log.Println(err)
		}
	}()
	return u, err
}

func (repo *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(u), nil
}

func (repo *CacheUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

func (repo *CacheUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
	}
}
