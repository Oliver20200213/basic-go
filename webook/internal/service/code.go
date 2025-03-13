package service

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/service/sms"
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
	"math/rand"
	"time"
)

// const codeTplId = "1877556"
// 改造成实时更新的配置
var codeTplId = atomic.String{}

var (
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
	//tplId string  模板id一般很少变动，可以直接定义，不用每次初始化的时候传入
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	codeTplId.Store("1877556") // 先初始化一个值
	// 启动onchange
	viper.OnConfigChange(func(in fsnotify.Event) {
		codeTplId.Store(viper.GetString("code.tpl.id"))
	})
	return &codeService{repo: repo, smsSvc: smsSvc}
}

// Send 发验证码，需要什么参数
func (svc *codeService) Send(ctx context.Context,
	biz string, // 区别是什么业务场景
	//code string, //这个码，谁来生成，一般可以自己生成
	phone string) error {
	// 三个步骤：生成一个验证码
	code := svc.generateCode()
	// 塞进去Redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		// 有问题
		return err
	}
	// 发送出去
	err = svc.smsSvc.Send(ctx, codeTplId.Load(), []string{code}, phone) //codeTplId.Load()获取当前codeTplId的值
	//if err != nil {
	//	// 这个地方怎么办？ 是否引入重试？
	//	// 这意味着，Redis 有这个验证码，但是短信没有发送成功，用户根本收不到
	//	// 我能不能删掉这个验证码？ 不能
	//	// 这里的这个err可能是超时的err， 你都不知道，发出了没有
	//	// 能不能重试？
	//	// 要重试的话，初始化的时候，传入一个自己就会重试 smsSvc
	//
	//}
	return err
}

// 区别出验证码对还是不对，或者是验证码业务有问题
func (svc *codeService) Verify(ctx context.Context, biz string,
	phone string, inputCode string) (bool, error) {
	//redis中key的形式：
	// phone_code:$biz:$phone
	// phone_code:login:152xxxxxx
	// code:login:152xxxxx
	// user:login:code:152xxxxx

	return svc.repo.Verify(ctx, biz, phone, inputCode)

}

//// 也可以值返回error,如果error为空则成功
//func (svc *codeService) VerifyV1(ctx context.Context, biz string,
//	phone string, inputCode string) error {
//
//}

func (svc *codeService) generateCode() string {
	// 六位数，num 在0,999999之间，包含0和999999
	num := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000)
	// 不够六位的加上前导0
	// 000001
	return fmt.Sprintf("%06d", num)
}
