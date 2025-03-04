package auth

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

/*
增加静态token认证，内部短信的调用
*/

type SMSService struct {
	svc sms.Service
	key string
}

func NewSMSService(svc sms.Service, key string) sms.Service {
	return &SMSService{
		svc: svc,
		key: key,
	}
}

// Send 发送 其中biz必须是线下申请的一个代表业务方的token
func (s *SMSService) Send(ctx context.Context, biz string,
	args []string, numbers ...string) error {

	var tc TokenClaims
	// 在这里进行权限校验
	// 如果我这里能解析成功，说明就是对应的业务方
	// 解析对了没有err 就说明token是我发的
	token, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token 不合法")
	}
	return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}

type TokenClaims struct {
	jwt.RegisteredClaims
	Tpl string
}

/*
回顾jwt相关的内容
jwt的解析
token,err := jwt.ParseWithClaims(biz,&tc,func(token *jwt.Token)(interface{},error){
	return s.key, nil
})
biz, 这是你要解析的jwt字符串
&tc，这是实现了jwt.Claims接口的结构体实例，用于存储jwt中的声明，通常是传递一个结构体的指针，以便于存储。
func(token,*jwt.Token)(interface{},error){},这是一个回调函数，用于获取验证jwt签名所需要的秘钥，
在解析jwt时，这个函数会被调用，你需要在这个函数中返回用验证签名的秘钥

token，是一个*jwt.Token对象，包含了jwt的各个部分（Header，Claims,Signature）



jwt token的生成
先生成jwt.Token对象
type TokenClaims struct{
	Tpl string
	jwt.RegisteredClaims
}
claims := TokenClaims{
	Tpl:"短信服务商提供的tpl"
	RegisteredClaims:jwt.RegisteredClaims{
		// 内部使用可以将token设置为不过期，方式1：不设置ExpiresAt,过期时间，方式2：设置一个非常遥远的过期时间
		ExpiresAt:jwt.NewNumericDate(time.Now().Add(time.Hour*24*365*100))
	}
}
token:= jwt.NewWithClaims(jwt.SigningMethodHS512,claims)
jwt.SigningMethodHS512：加密算法
claims：jwt.Claims的结构体实例

通过key生成jwt token的字符串
tokenStr, err :=token.SignedString([]byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"))
[]byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90")：加密的秘钥
*/
