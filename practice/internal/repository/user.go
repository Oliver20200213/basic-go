package repository

import (
	"basic-go/practice/internal/domain"
	"basic-go/practice/internal/repository/cache"
	"basic-go/practice/internal/repository/dao"
	"context"
	"log"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, cache *cache.UserCache) *UserRepository {
	return &UserRepository{dao: dao, cache: cache}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil

}

func (repo *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
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
	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}
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
