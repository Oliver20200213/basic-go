package web

import (
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/oauth2/wechat"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	JwtHandler
	stateKey []byte
	cfg      WechatHandlerConfig
}

type WechatHandlerConfig struct {
	Secure bool
	//stateKey []byte 也可以将stateKey放到这里
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService,
	cfg WechatHandlerConfig) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:      svc,
		userSvc:  userSvc,
		stateKey: []byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"),
		cfg:      cfg,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.POST("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登录URL失败",
		})
	}

	err = h.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
	}

	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间为你预期中一个用户完成登录的时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return err
	}

	// 详细说明见下面注释
	ctx.SetCookie("jwt-state", tokenStr, 600,
		"/oauth2/wechat/callback",
		"", h.cfg.Secure, true) // 正常线上要将secure设置为true
	return nil
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// 临时测试：修改自己本机的hosts，将回调的url指向自己的电脑
	//ctx.String(http.StatusOK , "你过来了！")

	code := ctx.Query("code")
	// 验证state
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登录失败",
		})
	}

	// 验证微信的code
	info, err := h.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误", // 严格的来说得区分一下err，是不是攻击者伪造的code
		})
		// 做好监控
		// 做好日志记录
		return
	}

	// 这里怎么办，设置jwt token（之前是在userhandler下面的需要摘出来）
	// uid从哪来？
	// 从userService里面拿uid
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

	err = h.setLoginToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})

}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	// 校验一下获取到的state
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		// 有人搞事
		// 做好监控
		// 做好日志记录
		return fmt.Errorf("拿不到state的cookie，%w", err)
	}

	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		// 做好监控
		// 做好日志记录
		return fmt.Errorf("token已经过期，%w", err)
	}
	if sc.State != state {
		// 做好监控
		// 做好日志记录
		return errors.New("state不相等")
	}
	return nil
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}

// 另一种设计思路，从url的路径参数中获取是使用的哪个登录，然后再分发
//type OAuth2Handler struct {
//	wechatService
//	dingdingService
//	feishuService
//
//	//或者
//	svcs map[string]OAuth2
//}
//
//func NewOAuth2Handler() {}
//func (h *OAuth2Handler) RegisterRoutes(server *gin.Engine) {
//	g := server.Group("/oauth2")
//	g.POST("/:platform/authurl", h.AuthURL)
//	g.Any("/:platform/callback", h.Callback)
//}
//
//func (h *OAuth2Handler) AuthURL(ctx *gin.Context) {
//	platform := ctx.Param("platform") // 获取路径参数
//	switch platform { // 根据路径参数分发
//	case "wechat":
//		h.wechatService.AuthURL
//	}
//
//	// map的写法
//	svc := h.svcs[platform]
//	svc.AuthUrl
//}

/*
接口：需要构造两个接口
接口1：用于构造跳转到微信那边的url
接口2：用于处理微信跳转回来的请求（也就是接口1中url中的redirectURI）

接口1：
/oauth2/wechat/authurl   前端点击微信登录时，访问该接口
后台根据微信开放平台的要求构造authurl返回给前端，const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
前端拿到这个url后访问，出现微信扫码页面

接口2：
/oauth2/wechat/callback 当用户扫码确认之后会跳转到该接口
去url中获取微信返回的code进行验证，使用该code构造url，	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
get访问后获取到access_code等信息，然后使用jwt token保存用户登录状态


openid和unionid的区别
openid在当前应用中唯一
unionid在公司中唯一
在某个应用内使用openid，夸应用使用unionid




"jwt-state":
含义: Cookie 的名称。
解释: 这个参数指定了 Cookie 的名称为 "jwt-state"。

tokenStr:
含义: Cookie 的值。
解释: 这个参数是 Cookie 的实际内容，通常是一个字符串。在这里，tokenStr 是一个变量，表示要存储在 Cookie 中的 JWT（JSON Web Token）或其他类型的令牌。

600:
含义: Cookie 的过期时间（以秒为单位）。
解释: 这个参数指定了 Cookie 的有效期为 600 秒（即 10 分钟）。超过这个时间后，Cookie 将自动过期。

"/oauth2/wechat/callback":
含义: Cookie 的路径。
解释: 这个参数指定了 Cookie 的有效路径。只有在这个路径下的请求才会携带这个 Cookie。在这里，路径是 "/oauth2/wechat/callback"，意味着只有在这个路径下的请求才会包含这个 Cookie。

"":
含义: Cookie 的域名。
解释: 这个参数指定了 Cookie 的有效域名。如果为空字符串，则表示当前域名。你可以将其设置为特定的域名，以限制 Cookie 只在那个域名下有效。

false:
含义: 是否仅通过 HTTPS 传输 Cookie。
解释: 这个参数指定了 Cookie 是否只能通过 HTTPS 协议传输。如果设置为 true，则 Cookie 只能在 HTTPS 连接中传输；如果设置为 false，则 Cookie 可以通过 HTTP 或 HTTPS 传输。

true:
含义: 是否禁止客户端 JavaScript 访问 Cookie。
解释: 这个参数指定了 Cookie 是否启用 HttpOnly 标志。如果设置为 true，则客户端 JavaScript 无法访问这个 Cookie，这有助于防止 XSS（跨站脚本攻击）攻击
*/
