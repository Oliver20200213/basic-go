package cache

import (
	"basic-go/webook/internal/repository/cache/redismocks"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "验证码可以发送",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				// 由于代码中c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code)
				// 返回的是*Cmd的类需要手动构造
				res := redis.NewCmd(context.Background())
				// 由于c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
				// 最后返回的是int类型的值，进行检测的也是int的值，需要手动构造 res.SetVal()
				//res.SetErr(nil) 默认
				res.SetVal(int64(0)) //验证码可以发送的值是0 这里只能接受int64或者string类型
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:15512345678"},
					[]any{"123456"}). // 由于Eval最有一个参数为不定参数所以这样写
					Return(res)       // 注意这里返回的是Eval的返回值*Cmd
				return cmd
			},
			biz:     "login",
			phone:   "15512345678",
			code:    "123456",
			wantErr: nil,
		},
		{
			name: "验证码发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				// 构造*Cmd
				res := redis.NewCmd(context.Background())
				// 由于最后检测的int类型，需要手动构造
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:15512345678"},
					[]any{"123456"}).
					Return(res)
				return cmd
			},
			biz:     "login",
			phone:   "15512345678",
			code:    "123456",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("mock redis error"))
				//res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:15512345678"},
					[]any{"123456"}).
					Return(res) // 注意这里返回的是Eval的返回值*Cmd
				return cmd
			},
			biz:     "login",
			phone:   "15512345678",
			code:    "123456",
			wantErr: errors.New("mock redis error"),
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				//res.SetErr(nil)
				res.SetVal(int64(-10))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:15512345678"},
					[]any{"123456"}).
					Return(res) // 注意这里返回的是Eval的返回值*Cmd
				return cmd
			},
			biz:     "login",
			phone:   "15512345678",
			code:    "123456",
			wantErr: errors.New("系统错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCodeCache(tc.mock(ctrl))
			err := c.Set(context.Background(), tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
