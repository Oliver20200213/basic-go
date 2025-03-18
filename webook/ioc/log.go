package ioc

import (
	"basic-go/webook/pkg/logger"
	"go.uber.org/zap"
)

func InitLogger() logger.LoggerV1 {
	// 全局定义的日志，如果想要自己单独使用的不想和别人共享的，可以用装饰器单独的包一层，例子见ioc中的user.go
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// zap.ReplaceGlobals(l) 满足不完全注入的zap.L()的使用
	return logger.NewZapLogger(l)
	// 直接返回了了Logger，那个模块使用可以通过依赖注入直接使用这个Logger
}

/*
小结：
初始化：
logger, err:= zap.NewDevelopment()
使用：
logger可以直接使用，或者通过zap.L()来使用
通过zap.L()使用的时候需要先用初始化好的logger替换掉默认的全局logger，zap.ReplaceGlobals(logger)
否则通过zap.L()不能打印出日志来

*/
