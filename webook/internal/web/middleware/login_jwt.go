package middleware

import (
	ijwt "basic-go/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

// JWT登录校验

type LoginJWTMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(jwtHandler ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHandler,
	}
}

// IgnorePaths 中间方法，用于构建部分
func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

// Build终结方法，返回你最终希望的数据
func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//不需要的校验
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		////使用JWT来校验
		//tokenHeader := ctx.GetHeader("Authorization")
		//if tokenHeader == "" {
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		////此时tokenHeader的格式：Bearer xdsfdsadsfd最前面是有个Bearer的
		//segs := strings.Split(tokenHeader, " ")
		////segs := strings.SplitN(tokenHeader, ".", 2) 或者直接切2次
		//if len(segs) != 2 {
		//	//没有登录,有人捣乱
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		// 改为长短token
		tokenStr := l.ExtractToken(ctx)
		claims := &ijwt.UserClaims{}
		//ParseWithClaims里面一定要传入claims指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"), nil
		})
		if err != nil {
			//没有登录，Bearer xxxxx
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//默认是不用自己手动验证过期时间的，会自己验证如果过期token.Valid会为false
		//如果实在想自己验证，也可以自己比较验证
		//if claims.ExpiresAt.Before(time.Now()){
		//	//过期了
		//}
		//err为nil，token不为nil
		if token == nil || !token.Valid || claims.Uid == 0 {
			//没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			//严重的安全问题
			//是需要加监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 这里可以考虑一下降级操作，如果redis崩了，可以不进行严格的校验是不是已经主动退出登录了，直接通过
		//if redis 崩了{
		//	return
		//}

		// 查看redis中是否存储有当前的ssid（记录已经退出的ssid）
		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			// 要么 redis 有问题，要么已经退出登录了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 能不能检测短token过期了，搞个新的？
		// 这就和自动刷新没区别了，功能上是可以的，
		// 但是这所以使用长短token，是因为我们认为短toke被频繁使用个，更加容易泄密，
		// 因此才需要有一个长token，这个长token只在登录和调用refresh_token时使用，
		// 所以不容易泄露

		// 使用长短token之后这个刷新机制用不上了
		////更新过期时间
		////每十秒钟刷新一次
		//now := time.Now()
		//if claims.ExpiresAt.Sub(now) < time.Second*50 {
		//	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//	//需要重新生成token
		//	tokenStr, err = token.SignedString([]byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"))
		//	if err != nil {
		//		//记录日志
		//		log.Println("jwt续约失败", err)
		//	}
		//	//重新返回给前端
		//	ctx.Header("x-jwt-token", tokenStr)
		//}

		// web框架中相同gin.Context： Set是向ctx中添加数据    Get是获取数据ctx.Get("claims“).(string)
		// 如果是通用go程序context.Context：context.WithValue("token","xxxxx")向context.Context中添加数据, context.Value("token")获取数据
		ctx.Set("claims", claims)
	}
}

/*
面试降级策略
-  在redis没有崩溃的时候，就会严格的执行ssid的校验，判定用户有没有主动退出登录
- 如果redis崩溃，就不会严格的执行ssid的校验

相比退出登录的少数情况，大多数已经登录的用户并不会应为redis不可用而全部无法通过登录校验
*/
