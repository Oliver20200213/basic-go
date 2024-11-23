package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		emailRegexPaatern    = `^\w+(-+\.\w+)*@\w+(-.\w+)*.\w+(-.\w+)*$`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPaatern, regexp.None)
	pwdExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: pwdExp,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	up := server.Group("/users")
	up.POST("/signup", u.SignUp)
	up.POST("/login", u.Login)
	up.POST("/edit", u.Edit)
	up.GET("/profile", u.Profile)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.ShouldBind(&req); err != nil {
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码输入不一致")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，并且包含数字，特殊字符")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
	fmt.Printf("%v\n", req)
}
func (u *UserHandler) Login(ctx *gin.Context)   {}
func (u *UserHandler) Edit(ctx *gin.Context)    {}
func (u *UserHandler) Profile(ctx *gin.Context) {}
