package ratelimit

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/pkg/ratelimit"
	"context"
	"fmt"
)

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
		return errLimited // 不到逼不得已不要对外公开
	}
	err = s.Service.Send(ctx, tpl, args, numbers...) //相当于将 sms.Service 的方法直接包含在了 RateLimitSMSServiceV1 里面，因此可以直接使用 s.Service.Send()
	// 你在这也可以加一些代码，新特性
	return err
}
