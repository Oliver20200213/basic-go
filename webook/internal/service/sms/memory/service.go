package memory

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"fmt"
)

type Service struct {
}

func (s *Service) SendV1(ctx context.Context, tpl string, args []sms.NamedArg, numbers ...string) error {
	//TODO implement me
	panic("implement me")
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Println("Send：", args)
	return nil
}
