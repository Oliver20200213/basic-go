package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtHandler struct {
}

// 为什么这里不用指针，为了userhandler和wechathandler中组合使用
func (h JwtHandler) setJWTToken(ctx *gin.Context, uId int64) error {
	//如何在JWT token中携带数据，比如要带userId
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			//配置过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uId,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	//token.SigningString() 不使用key直接生成token
	//更安全使用key生成token
	tokenStr, err := token.SignedString([]byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"))
	if err != nil {
		return err
	}

	//将token放到header中
	ctx.Header("x-jwt-token", tokenStr) //将token放到header中
	fmt.Println(tokenStr)
	return nil
}

type UserClaims struct {
	jwt.RegisteredClaims
	//声明自己要放放进token里面的数据
	Uid int64
	//自己可以随便加，但是最好不要加敏感数据例如password 权限之类的信息
	UserAgent string
}
