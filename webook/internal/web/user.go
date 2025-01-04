package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2" //引入新的正则库替代标准库 这样引入可以使用regexp调用而不是regexp2
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// UserHandler 定义在它上面定义跟用户有关的路由
type UserHandler struct {
	svc         *service.UserService
	smsSvc      *service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService, smsSvc *service.CodeService) *UserHandler {
	const (
		emailRegexPattern = "^\\w+(-+.\\w+)*@\\w+(-.\\w+)*.\\w+(-.\\w+)*$"
		//使用``不用进行转义
		//emailRegexPattern := `^\w+(-+.\w+)*@\w+(-.\w+)*.\w+(-.\w+)*$`

		//注意go标准的正则库不支持复杂的正则，需要引入额外的库github.com/dlclark/regexp2
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
	)
	//先预编译
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None) //regexp.None第二个参数可以随便填
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)

	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		smsSvc:      smsSvc,
	}
}

// 另一种分组的方式,将分组放到外面
func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
	ug.POST("/signup.lua", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup.lua", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.ProfileJWT)
	// PUT "/login/sms/code" 发送验证码
	// POST "/login/sms/code" 校验验证码
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	const biz = "login"
	err := u.smsSvc.Send(ctx, biz, req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
	}
	ctx.String(http.StatusOK, "发送成功")

}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	//bind方法会根据content-type来解析你的数据到req里面
	//解析错了，就会直接写回一个400错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	//校验邮箱
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		//记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	//校验密码
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}

	//调用service方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	//errors.Is(err,service.ErrUserDuplicateEmail)  这个是错误的最佳实践
	//判断err是否和ErrUserDuplicateEmail相等，如果相等.web返回true
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	fmt.Printf("%v\n", req)
	ctx.String(http.StatusOK, "注册成功  ")
	return
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	//使用jwt设置登录态
	////生成一个JWT token（不带数据的）
	//token := jwt.New(jwt.SigningMethodHS512)

	//如何在JWT token中携带数据，比如要带userId
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			//配置过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	//token.SigningString() 不使用key直接生成token
	//更安全使用key生成token
	tokenStr, err := token.SignedString([]byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	//将token放到header中
	ctx.Header("x-jwt-token", tokenStr) //将token放到header中
	fmt.Println(tokenStr)
	fmt.Println(user)
	ctx.String(http.StatusOK, "登录成功")
	return

}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//session实现步骤2
	//登录成功之后
	//设置session
	sess := sessions.Default(ctx)
	//可以随便设置值了
	//要放在session里面的值
	sess.Set("userId", user.Id)
	//配置options,option是用来控制cookie的
	sess.Options(sessions.Options{
		//生产环境中需要开启下面两项
		//Secure:   true, //表示cookie需要用https才能发送给客户端
		//HttpOnly: true, //表示通过JavaScript无法访问cookie
		//HttpOnly: true确保cookie不会被客户端脚本访问，而Secure: true则确保cookie只能通过安全的HTTPS连接传输
		MaxAge: 60,
	})
	sess.Save()

	ctx.String(http.StatusOK, "登录成功")
	return

}
func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1, //maxage配置成小于表示cookie过期,同时还表示redis中key和value的过期时间
	})
	sess.Save()
	ctx.String(http.StatusOK, "退出登录成功")
}
func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		NickName string `json:"nikName"`
		BirthDay string `json:"birth"`
		Describe string `json:"describe"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	id, ok := sessions.Default(ctx).Get("userId").(int64)
	if !ok {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	err := u.svc.Edit(ctx, domain.User{
		Id:       id,
		NickName: req.NickName,
		BirthDay: req.BirthDay,
		Describe: req.Describe,
	})
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	//c, ok := ctx.Get("claims")
	////你可以断定，必然有claims
	//if !ok {
	//	//可以考虑监控住这里
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}

	//或者忽略这里的ok，直接用下面的断言来检测，如果claims值是nil那么断言会是false
	c, _ := ctx.Get("claims")
	//你可以断定，必然有claims
	//if !ok {
	//	//可以考虑监控住这里
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	//ok代表是不是*UserClaims
	claims, ok := c.(*UserClaims)
	if !ok {
		//可以考虑监控住这里
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	println(claims.Uid)

	user, err := u.svc.Profile(ctx, claims.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "用户信息：%s,", user.Email)

}

func (u *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "这是profile页面")
	//println("xxxxxx")
	return
}

type UserClaims struct {
	jwt.RegisteredClaims
	//声明自己要放放进token里面的数据
	Uid int64
	//自己可以随便加，但是最好不要加敏感数据例如password 权限之类的信息
	UserAgent string
}
