package web

//已全部整合到了jwt包中了

//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt/v5"
//	"github.com/google/uuid"
//	"strings"
//	"time"
//)
//
//type JwtHandler struct {
//	// access_token key
//	atKey []byte
//	// refresh_token
//	rtKey []byte
//}
//
//func NewJwtHandler() JwtHandler {
//	return JwtHandler{
//		atKey: []byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"),
//		rtKey: []byte("QAonYNt3DpoEojWkzJruRYmigFjmfn99"),
//	}
//}
//
//// 为什么这里不用指针，为了userhandler和wechathandler中组合使用
//func (h JwtHandler) setJWTToken(ctx *gin.Context, uId int64, ssid string) error {
//	//如何在JWT token中携带数据，比如要带userId
//	claims := UserClaims{
//		RegisteredClaims: jwt.RegisteredClaims{
//			//配置过期时间
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
//		},
//		Uid:       uId,
//		Ssid:      ssid,
//		UserAgent: ctx.Request.UserAgent(),
//	}
//	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
//
//	//token.SigningString() 不使用key直接生成token
//	//更安全使用key生成token
//	tokenStr, err := token.SignedString([]byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"))
//	if err != nil {
//		return err
//	}
//
//	//将token放到header中
//	ctx.Header("x-jwt-token", tokenStr) //将token放到header中
//	fmt.Println(tokenStr)
//	return nil
//}
//
//// 长短token
//func (h JwtHandler) setRefreshToken(ctx *gin.Context, uId int64, ssid string) error {
//	//如何在JWT token中携带数据，比如要带userId
//	claims := RefreshClaims{
//		RegisteredClaims: jwt.RegisteredClaims{
//			//配置过期时间
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
//		},
//		Uid:  uId,
//		Ssid: ssid,
//	}
//	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
//
//	//token.SigningString() 不使用key直接生成token
//	//更安全使用key生成token
//	tokenStr, err := token.SignedString(h.rtKey)
//	if err != nil {
//		return err
//	}
//
//	//将token放到header中
//	ctx.Header("x-refresh-token", tokenStr) //将token放到header中
//	fmt.Println(tokenStr)
//	return nil
//}
//
//func ExtractToken(ctx *gin.Context) string {
//	//tokenHeader := ctx.Request.Header.Get("Authorization") // context标准库中的方式
//	tokenHeader := ctx.GetHeader("Authorization") // gin框架中获取的方式
//	segs := strings.Split(tokenHeader, " ")
//	if len(segs) != 2 {
//		return "" // 如果返回空字符串，下一部解析的时候会直接爆粗，所以不用返回error直接返回空字符串即可
//	}
//	return segs[1]
//}
//
//func (h JwtHandler) setLoginToken(ctx *gin.Context, uId int64) error {
//	ssid := uuid.New().String()
//	err := h.setJWTToken(ctx, uId, ssid)
//	if err != nil {
//		return err
//	}
//	err = h.setRefreshToken(ctx, uId, ssid)
//	return err
//}
//
//type RefreshClaims struct {
//	Uid  int64
//	Ssid string
//	jwt.RegisteredClaims
//}
//
//type UserClaims struct {
//	jwt.RegisteredClaims
//	//声明自己要放放进token里面的数据
//	Uid  int64
//	Ssid string
//	//自己可以随便加，但是最好不要加敏感数据例如password 权限之类的信息
//	UserAgent string
//}
