package ioc

import (
	"basic-go/practice/internal/service/sms"
	"basic-go/practice/internal/service/sms/memory"
	"basic-go/practice/internal/service/sms/tencent"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
)

func InitSms() sms.Service {
	return memory.NewService()
}

func initTencentSmsService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("没有找到环境变量：SMS_SECRET_ID")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("没有找到环境变量：SMS_SECRET_KEY")
	}
	client, err := tencentSms.NewClient(common.NewCredential(secretId, secretKey), "ap-beijing",
		profile.NewClientProfile())
	if err != nil {
		panic("初始化腾客户端错误")
	}
	return tencent.NewService(client, "appId", "signName")

}
