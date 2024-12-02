package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2" //引入新的正则库替代标准库 这样引入可以使用regexp调用而不是regexp2
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserHandler 定义在它上面定义跟用户有关的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
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
	}
}

// 另一种分组的方式,将分组放到外面
func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
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
	sess.Save()

	ctx.String(http.StatusOK, "登录成功")
	return

}

func (u *UserHandler) Edit(ctx *gin.Context) {}

func (u *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "这是profile页面")
	return
}
