package ioc

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/internal/service/sms/memory"
	"basic-go/webook/internal/service/sms/tencent"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
)

func InitSMSService() sms.Service {
	// 可以方便的换内存 还是换阿里腾讯登sms服务
	return memory.NewService()
}

func initTencentSMSService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("没有找到环境变量：SMS_SECRET_ID")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("没有找到环境变量：SMS_SECRET_KEY")
	}

	c, err := tencentSms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-beijing",
		profile.NewClientProfile())
	if err != nil {
		panic("初始化客户端错误")
	}
	return tencent.NewService(c, "appId", "signName")
}
