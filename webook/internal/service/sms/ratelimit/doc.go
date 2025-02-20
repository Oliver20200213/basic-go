package ratelimit

/*
使用装饰器的方式来添加短信服务的限流服务
装饰器模式：
	不改变原有实现而增加新特性的一种设计模式

示例：以短信限流为例
type RateLimitSMSService struct{
	svc sms.Service   // 我要装饰sms.Service那么必然有一个字段是这个接口
	limiter ratelimit.Limiter // 这是新增加的用来限流的限流器
}

// 由于实现了send方法，也就是实现了sms.Service借口，所以返回的可以是sms.Service(鸭子类型)
func NewRateLimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service{
	return &RateLimitSMSService{
		svc: svc,
		limiter: limiter,
	}
}
func(s *RateLimitSMSService)Send(ctx, tpl string, args []string, numbers ...string) error{
	// 这里增加代码，新特性
	limited, err := s.limiter.Limit(ctx, "sms:ratelimit")
	if err != nil {
		return fmt.Errorf("短信服务判断是否限流出现错误，%w", err)
	}
	if limited {
		return fmt.Errorf("短信服务触发限流")
	}

	err = s.svc.Send(ctx,tpl, args, numbers...)

	// 这里也可以增加代码，新特性
	return err
}

*/
