package cache

import (
	"basic-go/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

////不推荐实践
//type Cache interface {
//	GetUser(ctx context.Context, id int64) (domain.User, error)
//	//读取文章
//	GetArticle(ctx context.Context, aid int64)
//	//其他业务
//	//...
//}
//
//// 如果有统一内存调用，这样方式为最佳实践
//type CacheV1 interface {
//	//你中间件团队去做的
//	Get(ctx context.Context, key string) (any, error)
//}
//type UserCache struct {
//	cache CacheV1
//}
//
//func (u *UserCache) GetUser(ctx context.Context, id int64) (domain.User, error) {}

// 弱一点的不依赖中间件的实践

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	//传单机Redis可以
	//传cluster的Redis也可以
	client redis.Cmdable
	//过期时间,如果不在这里可以加到GetUser中
	expiration time.Duration
}

// NewUserCache
// 面向接口编程和依赖注入的三板斧：
// A用到了B，B一定是接口  => 这个是保证面向接口
// A用到了B，B一定是A的字段  =>  为了规避包变量、包方法，这个了都非常缺乏扩展性
// A用到了B，A绝对不初始化B，而是外面注入 => 保持依赖注入(DI,Dependency Injection)和依赖反转(IOC)
func NewUserCache(client redis.Cmdable) *UserCache {
	//func NewUserCache(client redis.Cmdable,expiration time.Duration) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
		//不写死的话需要传入，但是一定要传入上面字段的格式time.Duration
		//而不是传入不同格式的数据，再进行转换，就算是要转换也要在外面转换
	}
}

// GetUser
// 只要err为nil 就认为缓存里面有数据
// 如果没有数据，返回一个特定的error
// func (cache *UserCache) Get(ctx context.Context, id int64,expiration time.Duration) (domain.User, error) {
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	//数据不存在，err= redis.Nil
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	//fmt.Println("cache user", u)
	return u, err

}
func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u) //不到万不得已不要忽略任何错误 val,_:=json.Marshal(u) 好处是可以更好的定位错误
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *UserCache) key(id int64) string {
	//key的命名：
	// user:info:123
	// user_info_123
	// bumen_xiaozu_yewu_info_key
	return fmt.Sprintf("user:info:%d", id)
}
