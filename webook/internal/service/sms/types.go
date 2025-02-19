package sms

import (
	"context"
)

// Service
// 发送消息需要的参数:目标手机号码，appid,签名，模板，参数
// appId和签名是可以固定的，在初始化的时候指定即可
type Service interface {
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
	//建议使用SendV1
	//SendV1(ctx context.Context, tpl string, args []NamedArg, numbers ...string) error
	//调用者需要知道实现者需要知道是什么参数，是[]string还是map[string]string
	//SendV2(ctx context.Context, tpl string, args any, numbers ...string) error
}

// 第二种方式，优雅的方式
type NamedArg struct {
	Val  string
	Name string
}
