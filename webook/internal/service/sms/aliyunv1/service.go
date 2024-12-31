package aliyunv1

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/ecodeclub/ekit"
)

type Service struct {
	client *sms.Client
}

func NewService(client *sms.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) SendSms(ctx context.Context, signName, tplCode string, phone []string) error {
	phoneLen := len(phone)
	for i := 0; i < phoneLen; i++ {
		phoneSignle := phone[i]

		// 1.生成验证码
		// code代码有问题
		code := fmt.Sprintf("%x", md5.Sum([]byte(phoneSignle)))
		//完全没有做成一个独立的发短信的实现，而是一个强耦合验证码的实现
		bcode, err := json.Marshal(map[string]interface{}{
			"code": code,
		})
		if err != nil {
			return err
		}

		// 2.初始化短信结构体
		smsRequest := &sms.SendSmsRequest{
			SignName:      ekit.ToPtr[string](signName),
			TemplateCode:  ekit.ToPtr[string](tplCode),
			PhoneNumbers:  ekit.ToPtr[string](phoneSignle),
			TemplateParam: ekit.ToPtr[string](string(bcode)),
		}

		// 3.发送短信
		smsResponse, err := s.client.SendSms(smsRequest)
		if err != nil {
			return err
		}
		if *smsResponse.Body.Code == "OK" {
			fmt.Println(phoneSignle, string(bcode))
		}
		return errors.New(*smsResponse.Body.Message)
	}
	return nil
}
