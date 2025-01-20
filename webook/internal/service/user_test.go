package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	repomocks "basic-go/webook/internal/repository/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func Test_userService_Login(t *testing.T) {
	// 做成一个测试用例都用到的时间
	now := time.Now()
	testCase := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository
		//输入
		ctx      context.Context
		email    string
		password string

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功", // 用户名和密码是对得上的
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 由于context 没有用到，所以可以直接用gomock.Any()
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					// FindByEmail返回什么这里就返回什么
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$x26flcsBGgWbaOGw18nFQesFFWKewFkx8.6edRCDAMRdRv8OXnVmy",
						Phone:    "15512345678",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "password@world123",

			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$x26flcsBGgWbaOGw18nFQesFFWKewFkx8.6edRCDAMRdRv8OXnVmy",
				Phone:    "15512345678",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123@qq.com",
			password: "password@world123",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("mock db 错误"))
				return repo
			},
			email:    "123@qq.com",
			password: "password@world123",

			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$x26flcsBGgWbaOGw18nFQesFFWKewFkx8.6edRCDAMRdRv8OXnVmy",
						Phone:    "15512345678",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "11password@world123",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 具体测测试代码
			svc := NewUserService(tc.mock(ctrl))
			// 要测是的接口
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			// 或者在上面字段上不定义直接是context.Background
			// user, err := usersvc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantUser, user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

// 输出加密后的密码
func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("password@world123"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
