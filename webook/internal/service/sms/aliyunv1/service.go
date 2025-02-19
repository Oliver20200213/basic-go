package aliyunv1

import (
	mysms "basic-go/webook/internal/service/sms"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/ecodeclub/ekit"
	"strings"
)

type Service struct {
	client   *sms.Client
	signName string
}

func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) SendV1(ctx context.Context, tpl string, args []mysms.NamedArg, numbers ...string) error {
	argsMap := make(map[string]string, len(args))
	for _, arg := range args {
		argsMap[arg.Name] = arg.Val
	}
	code, err := json.Marshal(argsMap)
	if err != nil {
		return err
	}

	smsRequest := &sms.SendSmsRequest{
		TemplateCode:  ekit.ToPtr[string](tpl),
		SignName:      ekit.ToPtr[string](s.signName),
		PhoneNumbers:  ekit.ToPtr[string](strings.Join(numbers, ",")),
		TemplateParam: ekit.ToPtr[string](string(code)),
	}

	smsResponse, err := s.client.SendSms(smsRequest)
	if err != nil {
		return err
	}
	if *smsResponse.Body.Code != "OK" {
		return errors.New(fmt.Sprintf("发送短信失败,code:%s", *smsResponse.Body.Code))
	}
	return nil
}

func NewService(client *sms.Client, signName string) *Service {
	return &Service{
		client:   client,
		signName: signName,
	}
}

//func (s *Service) SendSms(ctx context.Context, signName, tplCode string, phone []string) error {
//	phoneLen := len(phone)
//	for i := 0; i < phoneLen; i++ {
//		phoneSignle := phone[i]
//
//		// 1.生成验证码
//		// code代码有问题
//		code := fmt.Sprintf("%x", md5.Sum([]byte(phoneSignle)))
//		//完全没有做成一个独立的发短信的实现，而是一个强耦合验证码的实现
//		bcode, err := json.Marshal(map[string]interface{}{
//			"code": code,
//		})
//		if err != nil {
//			return err
//		}
//
//		// 2.初始化短信结构体
//		smsRequest := &sms.SendSmsRequest{
//			SignName:      ekit.ToPtr[string](signName),
//			TemplateCode:  ekit.ToPtr[string](tplCode),
//			PhoneNumbers:  ekit.ToPtr[string](phoneSignle),
//			TemplateParam: ekit.ToPtr[string](string(bcode)),
//		}
//
//		// 3.发送短信
//		smsResponse, err := s.client.SendSms(smsRequest)
//		if err != nil {
//			return err
//		}
//		if *smsResponse.Body.Code == "OK" {
//			//发送成功！
//			fmt.Printf("发送手机号: %s 的短信成功,验证码为【%s】\n", phoneSignle, code)
//		}
//		return errors.New(*smsResponse.Body.Message)
//	}
//	return nil
//}
