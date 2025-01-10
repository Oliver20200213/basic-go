package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("我也不知道发生了什么，反正是跟code有关")
)

// 编译器会在编译的时候，把set_code的代码放进来这个luaSetCode变量里
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify-code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

// NewCodeCacheGoBestPractice GO的最佳实践是返回具体类型
func NewCodeCacheGoBestPractice(client redis.Cmdable) *RedisCodeCache {
	return &RedisCodeCache{client: client}
}

// NewCodeCache wire实现需要返回的是接口而不是具体类型
func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{client: client}
}

func (c *RedisCodeCache) Set(ctx context.Context,
	biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	// c.client.Eval(ctx, lua脚本，key,code).Int()
	//	这里使用Int()是因为lua脚本中会返回 0 -1 -2数字
	if err != nil {
		return err
	}
	switch res {
	case 0:
		//毫无问题
		return nil
	case -1:
		//发送太频繁
		return ErrCodeSendTooMany
	default:
		//系统错误
		return errors.New("系统错误")
	}

}

// 可以自定义key的格式
func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		// 正常来说，如果频繁出现这个错误，要进行告警，应为有人在搞你
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
	default:
		return false, ErrUnknownForCode
	}
}

//可以只返回error
//func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) error { }

// LocalCodeCache 假如要切换到本地内存 就需要把lua脚本的逻辑在这里在写一遍
type LocalCodeCache struct {
}
