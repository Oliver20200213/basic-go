package ratelimit

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/pkg/ratelimit"
	"context"
	"fmt"
)

/*
使用组合的方式实现装饰器
两种方法的对比：
  - 使用组合
    用户可以直接访问sms.Service,绕开你装饰器本身
    可以只实现sms.Service的部分方法
  - 不使用组合：
    可以有效的组织用户绕开装饰器
    必须实现Service的全部方法

如果sms.Service中有很多种方法但是你就只需要装饰其中的一种或者是几种那么使用组合的方式
*/
type RateLimitSMSServiceV1 struct {
	sms.Service // 使用了匿名字段也就是可以直接调用sms.Service中的方法
	limiter     ratelimit.Limiter
}

func NewRateLimitSMSServiceV1(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RateLimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

// 装饰器实现限流
func (s *RateLimitSMSServiceV1) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	// 你这里可以加一些代码，新特性
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		// 系统错误
		// 可以限流：保守策略，你的下游很坑的时候，
		// 可以不限流：你的下游很强，业务可用性要求很高，尽量容错策略
		// 包一下这个错误
		return fmt.Errorf("短信服务判断是否限流出现问题，%w", err)

	}
	if limited {
		return errLimited // 不到逼不得已不要对外暴漏出去
	}
	err = s.Service.Send(ctx, tpl, args, numbers...) //相当于将 sms.Service 的方法直接包含在了 RateLimitSMSServiceV1 里面，因此可以直接使用 s.Service.Send()
	// 你在这也可以加一些代码，新特性
	return err
}

/*
总结：
开闭原则 非侵入式 装饰器 这三个经常一起出现

开闭原则：
	对修改闭合，对扩展开放
非侵入式：
	不修改有代码

记住一句话：
	侵入式修改是万恶之源。它会降低代码可读性，降低可测试性，强耦合，降低可扩展性。
侵入式 = 垃圾代码，这个等式基本成立。
除非逼不得已，不然绝对不要搞侵入式修改！！！
*/
