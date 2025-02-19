package tencent

import (
	mysms "basic-go/webook/internal/service/sms"
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
	//limiter  ratelimit.Limiter
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
		client:   client,
	}
}

//func NewService(client *sms.Client, appId string, signName string, limiter ratelimit.Limiter) *Service {
//	return &Service{
//		appId:    ekit.ToPtr[string](appId),
//		signName: ekit.ToPtr[string](signName),
//		client:   client,
//		limiter:  limiter,
//	}
//}

//TemplateParam的参数args格式:
//腾讯云的参数args是 []*string
//阿里云的参数args是 string， json串

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//// 侵入式的写法（改了已有的代码，不推荐）需要用装饰器实现
	//limited, err := s.limiter.Limit(ctx, "sms:tencent")
	//if err != nil {
	//	// 系统错误
	//	// 可以限流：保守策略，你的下游很坑的时候，
	//	// 可以不限流：你的下游很强，业务可用性要求很高，尽量容错策略
	//	// 包一下这个错误
	//	return fmt.Errorf("短信服务判断是否限流出现问题，%w", err)
	//
	//}
	//if limited {
	//	return fmt.Errorf("短信服务触发了限流")
	//}

	//如果是微服务有别人调用需要检测下，要发送的号码是不是空的
	//if len(numbers) == 0 {
	//	return errors.New("no numbers provided")
	//}
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId          //短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId
	req.SignName = s.signName          //短信签名
	req.TemplateId = ekit.ToPtr(tplId) //短信模板的id
	//需要将numbers转成切片
	// 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	req.PhoneNumberSet = s.toStringPtrSlice(numbers) //下发手机号码
	req.TemplateParamSet = s.toStringPtrSlice(args)  // 用于设置短信模板中的变量参数,
	// 例如：短信模板是：验证码是：${code}，${code}是一个变量，表示验证码
	// args = []string{"1234",} 这里${code}就是1234
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败 %s, %s ", *status.Code, *status.Message)
		}
	}
	return nil
}

// 加入阿里云短信之后，修改types.go中args类型跟随着更改
func (s *Service) SendV1(ctx context.Context, tplId string, args []mysms.NamedArg, numbers ...string) error {

	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](tplId)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = slice.Map[mysms.NamedArg, *string](args, func(idx int, src mysms.NamedArg) *string {
		return &src.Val
	})
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败 %s, %s ", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}

/*
利用泛型转换slice中元素类型的方法Map
slice.Map函数：会将slice中元素装换成其他类型并返回新的切片
func Map[Src any, Dst any](src []Src, m func(idx int, src Src) Dst) []Dst {
	dst := make([]Dst, len(src))
	for i, s := range src {
		dst[i] = m(i, s)
	}
	return dst
}
*/
