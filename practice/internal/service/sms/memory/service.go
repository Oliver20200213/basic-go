package memory

import (
	"basic-go/practice/internal/service/sms"
	"context"
	"fmt"
)

type Service struct{}

func NewService() sms.Service {
	return &Service{}
}

func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Println("内存模拟发送短信：", tpl, args, numbers)
	return nil
}
