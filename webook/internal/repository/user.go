package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
	"log"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFind        = dao.ErrUserNotFund
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

	//create的时候可以用下缓存，来提高性能，一般注册之后都会登录
	//单独在login的时候没有必要使用缓存
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
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
	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}
	// 将数据写入到cache中，可以开个goroutine
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 这里该怎么办？
			// 打日志做监控即可
			log.Println(err)
		}
	}()
	return u, err
}

// 注意：使用缓存的两个核心问题：第一个是一致性问题， 第二个是我的缓存崩了
