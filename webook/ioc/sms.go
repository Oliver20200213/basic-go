package ioc

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	// 可以方便的换内存 还是换阿里腾讯登sms服务
	return memory.NewService()
}
