package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		appId:    appId,    //SDKAppID,应用id，需要到腾云与短信控制台里面的应用列表中查找
		signName: signName, //短信息签名
		client:   client,
	}
}

func (svc *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = svc.appId                        //短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId
	req.SignName = svc.signName                        //短信签名
	req.TemplateId = ekit.ToPtr(tplId)                 //短信模板的id
	req.PhoneNumberSet = svc.toStringPtrSlice(numbers) //下发手机号码
	req.TemplateParamSet = svc.toStringPtrSlice(args)  // 用于设置短信模板中的变量参数,
	// 例如：短信模板是：验证码是：${code}，${code}是一个变量，表示验证码
	// args = []string{"1234",} 这里${code}就是1234

	resp, err := svc.client.SendSms(req)
	if err != nil {
		return err
	}

	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *status.Code != "ok" {
			return fmt.Errorf("发送短信失败：%s, %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (svc *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
