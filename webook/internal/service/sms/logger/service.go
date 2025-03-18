package logger

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"go.uber.org/zap"
)

type Service struct {
	svc sms.Service
}

// Send 通过装饰器来统一打debug日志
func (s Service) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	zap.L().Debug("发送短信", zap.String("biz", biz),
		zap.Any("args", args))
	err := s.svc.Send(ctx, biz, args, numbers...)
	if err != nil {
		zap.L().Debug("发送短信出现异常", zap.Error(err))
	}

	return err
}
