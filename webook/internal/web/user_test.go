package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	svcmocks "basic-go/webook/internal/service/mocks"
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncrypt(t *testing.T) {
	password := "hello#world123"
	//加密
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	//比较
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)
}

func TestNil(t *testing.T) {
	testTypeAssert(nil)
}

func testTypeAssert(c any) {
	claims := c.(*UserClaims)
	println(claims.Uid)
}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(nil)
				//usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				// 注册成功是return nil的
				return usersvc
			},
			reqBody: `
{
	"email":"123@qq.com",
	"password":"hello#world123",
	"confirmPassword":"hello#world123"
}
`,
			wantCode: 200,
			wantBody: "注册成功",
		},
		{
			name: "参数不对，bind失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email":"123@qq.com",
	"password":"hello#world123"
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email":"123com",
	"password":"hello#world123",
	"confirmPassword":"hello#world123"
}
`,
			wantCode: 200,
			wantBody: "你的邮箱格式不对",
		},
		{
			name: "两次输入密码不匹配",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email":"123@qq.com",
	"password":"helloworld123",
	"confirmPassword":"hello#world123"
}
`,
			wantCode: 200,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email":"123@qq.com",
	"password":"hello123",
	"confirmPassword":"hello123"
}
`,
			wantCode: 200,
			wantBody: "密码必须大于8位，包含数字、特殊字符",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(service.ErrUserDuplicateEmail)
				// 注册成功是return nil的
				return usersvc
			},
			reqBody: `
{
	"email":"123@qq.com",
	"password":"hello#world123",
	"confirmPassword":"hello#world123"
}
`,
			wantCode: 200,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(errors.New("随便返回一个错误"))
				// 注册成功是return nil的
				return usersvc
			},
			reqBody: `
{
	"email":"123@qq.com",
	"password":"hello#world123",
	"confirmPassword":"hello#world123"
}
`,
			wantCode: 200,
			wantBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 0.初始化ctrl
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 思路：构造http请求，获取到http的响应
			// 1.准备一个gin.Engine.并注册路由
			server := gin.Default()
			// codesvc在signup中没有用到，可以传nil
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			// 2.准备请求
			req, err := http.NewRequest(http.MethodPost,
				"/users/signup",
				bytes.NewBuffer([]byte(tc.reqBody)))
			// 数据是json格式
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err) //由于是自己定义的 所以任务是一定没有错误，如果有错误直接终止测试
			// 这里就可以继续使用req

			// 3.准备接收相应的Recorder
			resp := httptest.NewRecorder() // ResponseRecorder实现了responseWriter接口存储了HTTP的响应
			//resp.Code http的响应码
			//resp.Header()  http响应的header
			//resp.Body http响应的内容

			// 4.这是HTTP请求进到Gin框架的入口
			// 当你这样调用的时候，Gin就会处理这个请求
			// 响应写回到resp里
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code) //断言判断resp.Code是否和wantCode相同
			assert.Equal(t, tc.wantBody, resp.Body.String())

		})
	}
}

/*
mock工具
github地址：https://github.com/uber-go/mock
分成两个部分：
mockgen:命令行工具
测试中使用改的空mock对象的包
安装命令行工具：
go install go.uber.org/mock/mockgen@latest

windows下需要用绝对路径：
mockgen -source=E:/gowork/src/basic-go/webook/internal/service/user.go -package=svcmocks -destination=E:/gowork/src/basic-go/webook/internal/service/mocks/user.mock.go

*/
// mock的使用
func TestMock(t *testing.T) {
	// 先创建一个mock的控制器
	ctrl := gomock.NewController(t)
	// 每个测试结束都要调用Finish
	// 然后mock救护已验证你的测试流程是否符合预期
	defer ctrl.Finish()
	// svcmocks就是mockgen中定义的-package=svcmocks
	usersvc := svcmocks.NewMockUserService(ctrl)

	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
		//Times(2). 调用该方法两次
		Return(errors.New("mock error"))

	// 很容易犯的错误1，没有返回error（返回什么是看SignUp返回的是什么）
	//usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
	//	Return(123)

	// 易犯错误2：参数没对应
	//usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
	//	Email: "124@qq.com", //预期输入是124@qq.com 但是下面输入的是123@qq.com
	//}).
	//	Return(errors.New("mock error"))

	err := usersvc.SignUp(context.Background(), domain.User{ //context.Background()是创建一个空白的上下文
		Email: "123@qq.com",
	})
	t.Log(err) //输出的就是上面Return mock error

	// 设计了几次SignUp的调用就只能调几次，多了少了都不行，并且顺序也不能错
	err = usersvc.SignUp(context.Background(), domain.User{ //context.Background()是创建一个空白的上下文
		Email: "123@qq.com",
	})
	t.Log(err)
}
