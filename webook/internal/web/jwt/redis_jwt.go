package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var (
	AtKey = []byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90")
	RtKey = []byte("QAonYNt3DpoEojWkzJruRYmigFjmfn99")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}

func (r *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	fmt.Println(tokenStr)
	return nil
}

func (r *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RtKey)

	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil

}

func (r *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func (r *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := r.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = r.SetRefreshToken(ctx, uid, ssid)
	return err
}

func (r *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	// 将长短token设置成无效值
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	c, _ := ctx.Get("claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		return errors.New("获取token失败")
	}
	// 将过期的token的ssid 写入redis
	return r.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid), "", time.Hour*7*24).Err()

}

func (r *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	cnt, err := r.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	if err != nil || cnt > 0 {
		// 要么redis有问题，要么已经退出登录了
		return errors.New("redis error or ssid已存在")
	}
	return nil
}
