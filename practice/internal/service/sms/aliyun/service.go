package aliyun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/ecodeclub/ekit"
	"math/rand"
	"time"
)

type Service struct {
	client *sms.Client
}

func NewService(client *sms.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Send(ctx context.Context, signName, tplId string, phone []string) error {
	for i := 0; i < len(phone); i++ {
		phoneNumber := phone[i]

		//随机生成6位数的随机数
		code := fmt.Sprintf("%06v",
			rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
		bcode, err := json.Marshal(map[string]string{
			"code": code,
		})
		if err != nil {
			return err
		}

		smsRequest := &sms.SendSmsRequest{
			SignName:      ekit.ToPtr(signName),
			TemplateCode:  ekit.ToPtr(tplId),
			PhoneNumbers:  ekit.ToPtr(phoneNumber),
			TemplateParam: ekit.ToPtr(string(bcode)),
		}

		smsResponse, err := s.client.SendSms(smsRequest)
		if err != nil {
			return err
		}
		if *smsResponse.Body.Code == "OK" {
			//发送成功！
			fmt.Printf("发送手机号: %s 的短信成功,验证码为【%s】\n", phoneNumber, code)
		}
		return errors.New(*smsResponse.Body.Message)

	}
	return nil
}
