package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2" //引入新的正则库替代标准库 这样引入可以使用regexp调用而不是regexp2
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const biz = "login"

// 确保 UserHandler上实现了handler接口
var _ handler = &UserHandler{}

// 第二种写法，这种写法更优雅
var _ handler = (*UserHandler)(nil)

// UserHandler 定义在它上面定义跟用户有关的路由
type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	JwtHandler
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	const (
		emailRegexPattern = "^\\w+(-+.\\w+)*@\\w+(-.\\w+)*.\\w+(-.\\w+)*$"
		//使用``不用进行转义
		//emailRegexPattern = `^\w+(-+.\w+)*@\w+(-.\w+)*.\w+(-.\w+)*$`

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
		codeSvc:     codeSvc,
		JwtHandler:  NewJwtHandler(),
	}
}

// 另一种分组的方式,将分组放到外面
//func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
//	ug.POST("/signup.lua", u.SignUp)
//	//ug.POST("/login", u.Login)
//	ug.POST("/login", u.Login)
//	ug.POST("/edit", u.Edit)
//	ug.GET("/profile", u.Profile)
//}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	//ug.GET("/profile", u.ProfileJWT)
	ug.GET("/profile", u.ProfileJWTV1)

	// PUT "/login/sms/code" 发送验证码
	// POST "/login/sms/code" 校验验证码
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSms)
	ug.POST("/refresh_token", u.RefreshToken)
}

// RefreshToken 可以同时刷新长短token， 用 redis 来记录是否有效，即 refresh_token是一次性的
// 也可以参考登录校验部分，比较User-Agent来增强安全性
func (u *UserHandler) RefreshToken(ctx *gin.Context) {
	// 如果在调用RefreshToken的时候超时了怎么办
	// 也就是旧的access_token已经无效，新的没有获得
	// 没有办法 只能重新登录

	// 正常访问的时候Authorization里面应该是短token,access_token
	// 当访问该接口的时候Authorization里面应该是长token,refresh_token
	// 使用refresh_token来刷新access_token
	refreshToken := ExtractToken(ctx)
	var rc RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return u.rtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 搞个新的access_token
	err = u.setJWTToken(ctx, rc.Uid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "access token 刷新成功",
	})

}

func (u *UserHandler) LoginSms(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 校验：是不是一个合法的手机号码
	// 考虑用正则表达式
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		return
	}
	// 在这之前，可以加上各种校验
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}
	// 这个手机号会不是一个新用户呢？
	// 这里如果没有注册过，同时注册
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 用户id要从那里获取
	if err = u.setJWTToken(ctx, user.Id); err != nil {
		// 记录日志
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if err = u.setRefreshToken(ctx, user.Id); err != nil {
		// 需要记录日志
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "验证码校验通过",
	})

}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	// 拿到手机号码
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 校验：是不是一个合法的手机号码
	// 考虑用正则表达式
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		return
	}

	// 发送验证码
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

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
		ctx.String(http.StatusOK, "系统错误")
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
	ctx.String(http.StatusOK, "注册成功")
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

	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if err = u.setRefreshToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
	fmt.Println(user)
	return

}

// 摘出去，可以让wechat登录能使用
//func (u *UserHandler) setJWTToken(ctx *gin.Context, uId int64) error {
//	//如何在JWT token中携带数据，比如要带userId
//	claims := UserClaims{
//		RegisteredClaims: jwt.RegisteredClaims{
//			//配置过期时间
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
//		},
//		Uid:       uId,
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
	type Profile struct {
		Email string
	}
	sess := sessions.Default(ctx)
	id := sess.Get("userIdKey").(int64) // 如果不是int64在断言的时候可能会panic掉
	user, err := u.svc.Profile(ctx, id)
	if err != nil {
		// 按道理来说，id是对应的数据是肯定存在的，所以要是没找到，
		// 那说明是系统出问题了
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: Profile{
			Email: user.Email,
		},
	},
	)
}

// 将UserClaims也摘到jwt.go中
//type UserClaims struct {
//	jwt.RegisteredClaims
//	//声明自己要放放进token里面的数据
//	Uid int64
//	//自己可以随便加，但是最好不要加敏感数据例如password 权限之类的信息
//	UserAgent string
//}

func (u *UserHandler) ProfileJWTV1(ctx *gin.Context) {
	type Profile struct {
		Email    string
		Phone    string
		Nickname string
		Birthday string
		AboutMe  string
	}
	uc := ctx.MustGet("claims").(*UserClaims) //与Get相比如果没有获取到user则会panic
	user, err := u.svc.Profile(ctx, uc.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 这里不建议将领域的domain.user暴漏出去，一个是domain.user中有可能会有敏感数据像密码啥啥的
	// 在一个是不知都以后你的同事会不会将其他敏感的信息添加到domain.user中，所以需要自己定义要返回的数据
	ctx.JSON(http.StatusOK, Result{
		Data: Profile{
			Email:    user.Email,
			Phone:    user.Phone,
			Nickname: user.Nickname,
			Birthday: user.Birthday.Format(time.DateOnly),
			// time.DateOnly是常量"2006-01-02" 这里是将time.Time主换成time.DateOnly格式的字符串
			AboutMe: user.AboutMe,
		},
	})

}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		// 注意：其他字段，尤其是密码、邮箱和手机
		// 修改的时候需要通过别的手段
		// 邮箱、手机、密码都需要验证
		Nickname string `json:"nickname"`
		// 2024-01-13
		Birthday string `json:"birthday"`
		AboutMe  string `json:"about_me"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	// 可以在这两进行校验
	// 例如要求Nickname不需不为空
	// 校验的规则取决于产品经理
	if req.Nickname == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "昵称不能为空",
		})
		return
	}
	if len(req.AboutMe) > 1024 {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "关于我过长",
		})
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday) // 将字符串是time.DateOnly格式（yy-mm-dd）解析为time.Time格式
	if err != nil {
		// 这里其实没有直接校验具体的格式
		// 如果能欧成功转化过来，那就寿命没问题
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "日期格式不对",
		})
		return
	}

	uc := ctx.MustGet("claims").(*UserClaims)
	err = u.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       uc.Uid,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})

}
