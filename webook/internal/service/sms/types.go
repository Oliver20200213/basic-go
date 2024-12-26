package sms

import "context"

// Service
// 发送消息需要的参数:目标手机号码，appid,签名，模板，参数
// appId和签名是可以固定的，在初始化的时候指定即可
type Service interface {
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
}
