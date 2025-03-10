package cache

import (
	"context"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"sync"
	"time"
)

// 技术选型考虑的点
// 1. 功能性：功能是否能够完全覆盖你的需求
// 2. 社区和支持度：社区是否活跃，文档是否齐全
//	以及百度(搜索引擎)能不能搜到你需要的各种信息，有没有帮你踩坑
// 3. 非功能性：
//	易用性(用户友好度)，学习曲线要平滑
//	扩展性(如果开源软件的某些功能需要定制，框架是否支持定制，以及定制的难度高不高)
//  性能(追求性能的功能，旺旺有能力自研 )

// LocalCodeCache 本地缓存实现
type LocalCodeCache struct {
	cache *lru.Cache
	lock  sync.Mutex // 普通锁，或者说写锁（互斥的只能一个人拿到锁）
	// 读写锁 rwLock sync.RWMutex可以多个人加读锁
	rwLock     sync.RWMutex
	expiration time.Duration
}

func NewLocalCodeCache(c *lru.Cache, expiration time.Duration) *LocalCodeCache {
	return &LocalCodeCache{
		cache:      c,
		expiration: expiration,
	}
}

//// 通过biz+phone来加锁
//func (l *LocalCodeCache) lockKey(biz, phone string) string {
//	return fmt.Sprintf("%s:%s", biz, phone)
//}
//
//// 获取锁(为biz+phone创建锁)
//func (l *LocalCodeCache) getLock(key string) *sync.Mutex {
//	lock, _ := l.lockMap.LoadOrStore(key, &sync.Mutex{})
//	return lock.(*sync.Mutex)
//}

func (l *LocalCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	// 加了读锁
	l.rwLock.RLock()
	// 释放了读锁
	l.rwLock.RUnlock()
	// 加了写锁
	l.rwLock.Lock()
	// 释放了写锁
	l.rwLock.Unlock()

	l.lock.Lock()
	defer l.lock.Unlock()
	// 这里可以考虑读写锁来优化，但是效果不会很好（通过biz+phone来加锁）
	// 因为你可以预期，大部分时候是要走到写锁里面的
	//key := l.key(biz, phone)
	//lock := l.getLock(key) // 获取对应 key 的锁
	//
	//lock.Lock()  // 加锁
	//defer lock.Unlock() // 解锁
	//now := time.Now()
	//val, ok := l.cache.Load(key) 	// 尝试获取缓存中的数据

	// 我选用的本地缓存，很不幸的是，没有获得过期时间的接口，所以
	key := l.key(biz, phone)
	now := time.Now()
	val, ok := l.cache.Get(key)
	if !ok {
		//说明没有验证码
		l.cache.Add(key, codeItem{
			code:   code,
			cnt:    3,
			expire: now.Add(l.expiration),
		})
		return nil
	}
	// 解析数据
	item, ok := val.(codeItem)
	if !ok {
		// 理论上来说这是不可能的,也不要忽略它
		return errors.New("系统错误")
	}
	if item.expire.Sub(now) > time.Minute*9 {
		// 发送间隔不到一分钟
		return ErrCodeSendTooMany
	}

	// 如果是到了1分钟就可以重发
	l.cache.Add(key, codeItem{
		code:   code,
		cnt:    3,
		expire: now.Add(l.expiration),
	})
	return nil
}

func (l *LocalCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("%s:%s", biz, phone)
}

type codeItem struct {
	code   string
	cnt    int
	expire time.Time
}
