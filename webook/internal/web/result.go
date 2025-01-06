package web

type Result struct {
	// 业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// 501001 => 表示验证码没有发送出去
// 5 系统错误
// 01 表示登录模块
// 001 表示在模块中具体的错误
