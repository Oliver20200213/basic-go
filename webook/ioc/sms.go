package ioc

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/internal/service/sms/memory"
	sratelimit "basic-go/webook/internal/service/sms/ratelimit"
	"basic-go/webook/internal/service/sms/retryable"
	"basic-go/webook/internal/service/sms/tencent"
	"basic-go/webook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	"time"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	// 可以方便的换内存 还是换阿里腾讯登sms服务
	//return memory.NewService()

	// ratelimit的使用
	svc := sratelimit.NewRateLimitSMSService(memory.NewService(), ratelimit.NewRedisSlidingWindowLimiter(cmd, time.Second, 3000))
	// 增加重试功能
	return retryable.NewService(svc, 3)

	//return memory.NewService()

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
