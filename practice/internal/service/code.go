package service

import (
	"basic-go/practice/internal/repository"
	"basic-go/practice/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
	"time"
)

var (
	ErrCodeSendToMany = repository.ErrCodeSendTooMany
)

const codeTpl = "2343627"

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *codeService) Send(ctx context.Context, biz, phone string) error {
	// 生成随机code
	code := svc.generateCode()
	// 检测能否发送
	err := svc.repo.Store(ctx, code, biz, phone)
	if err != nil {
		return err
	}
	// 发送验证码
	err = svc.smsSvc.Send(ctx, codeTpl, []string{code, "15"}, phone)
	return err
}

func (svc *codeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *codeService) generateCode() string {
	// 随机生成6位数
	num := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000)
	return fmt.Sprintf("%06v", num)
}
