package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	cachemocks "basic-go/webook/internal/repository/cache/mocks"
	"basic-go/webook/internal/repository/dao"
	daomocks "basic-go/webook/internal/repository/dao/mocks"
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	// 111ms.11111ns
	now := time.Now().In(time.Local)
	// 需要去掉毫秒以外的部分 111ms
	now = time.UnixMilli(now.UnixMilli())
	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		ctx      context.Context
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				// 缓存未命中，查了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)
				d := daomocks.NewMockUserDao(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{
						Id: 123,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "this is a password", // 密码可以随便填，因为repository.FindById中不涉及到密码的加密和解密
						Phone: sql.NullString{
							String: "15512345678",
							Valid:  true,
						},
						Ctime:    now.UnixMilli(),
						Utime:    now.UnixMilli(),
						Birthday: now.UnixMilli(),
					}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "this is a password",
					Phone:    "15512345678",
					Ctime:    now, //注意这里的时间需要保留到毫秒，time.Now()是到纳秒的
					Birthday: now,
				}).Return(nil)
				return d, c
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "this is a password",
				Phone:    "15512345678",
				Ctime:    now,
				Birthday: now,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				// 缓存命中
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDao(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{
						Id:       123,
						Email:    "123@qq.com",
						Password: "this is a password",
						Phone:    "15512345678",
						Ctime:    now,
						Birthday: now,
					}, nil)
				return d, c
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "this is a password",
				Phone:    "15512345678",
				Ctime:    now,
				Birthday: now,
			},
			wantErr: nil,
		},
		{
			name: "缓存未命中，查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				// 缓存未命中，查了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)
				d := daomocks.NewMockUserDao(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{}, errors.New("mock DB错误"))
				return d, c
			},
			ctx:      context.Background(),
			id:       123,
			wantUser: domain.User{},
			wantErr:  errors.New("mock DB错误"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			u, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantUser, u)
			assert.Equal(t, tc.wantErr, err)
			// 测试goroutine：等待1s确保能运行到goroutine中
			time.Sleep(time.Second)
		})
	}
}
