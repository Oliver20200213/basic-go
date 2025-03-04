package web

import (
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	JwtHandler
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.POST("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	url, err := h.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登录URL失败",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// 临时测试：修改自己本机的hosts，将回调的url指向自己的电脑
	//ctx.String(http.StatusOK , "你过来了！")

	// 验证微信的code
	code := ctx.Query("code")
	state := ctx.Query("state")
	info, err := h.svc.VerifyCode(ctx, code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误", // 严格的来说得区分一下err，是不是攻击者伪造的code
		})
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

	err = h.setJWTToken(ctx, u.Id)
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
*/
