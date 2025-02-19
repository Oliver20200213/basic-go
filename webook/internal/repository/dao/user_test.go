package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

/*
核心：
mock数据库：使用sqlmock来创建gorm的db，关闭乱七八糟的语句
*/
func TestGORMUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name string
		// 为什么不用ctrl？
		// 因为这里是sqlmock 不是gomock
		mock    func(t *testing.T) *sql.DB
		ctx     context.Context
		user    User
		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				// mockDB 是一个*sql.DB 对象，可以像操作真实数据库一样操作它
				// mock 是一个 sqlmock.Sqlmock 对象，用于设置模拟的 SQL 查询行为（如预期执行的 SQL 语句、返回的结果或错误）。
				mockDB, mock, err := sqlmock.New()
				// 构建insert之后的返回值
				res := sqlmock.NewResult(3, 1) // 3表示插入的id 1表示受影响的行数
				// 增删改用ExpectExec  查询用ExpectQuery
				// 后面是.*，这里预期的是一个正则表达式
				// 这个写法的意思是，只要是INSERT 到users的语句就行
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnResult(res)
				require.NoError(t, err)

				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				// 增删改用ExpectExec  查询用ExpectQuery
				// 后面是.*，这里预期的是一个正则表达式
				// 这个写法的意思是，只要是INSERT 到users的语句就行
				mock.ExpectExec("INSERT INTO `users`.*").
					WillReturnError(&mysql.MySQLError{
						Number: 1062,
					})
				require.NoError(t, err)

				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
			wantErr: ErrUserDuplicate,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				// 增删改用ExpectExec  查询用ExpectQuery
				// 后面是.*，这里预期的是一个正则表达式
				// 这个写法的意思是，只要是INSERT 到users的语句就行
				mock.ExpectExec("INSERT INTO `users`.*").
					WillReturnError(errors.New("db error"))
				require.NoError(t, err)

				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
			wantErr: errors.New("db error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			db, err := gorm.Open(gormMysql.New(gormMysql.Config{
				//Conn: mockDB,
				Conn: tc.mock(t),
				// 关闭 SELECT VERSION; 的调用(如果为false那么GORM在初始化的时候会先调用show version)
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// mock DB 不需要ping
				DisableAutomaticPing: true,
				// 关闭默认commit
				SkipDefaultTransaction: true,
			})
			d := NewUserDAO(db)
			err = d.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)

		})
	}
}

/*
sqlmock是和go gdk中database/mysql同级的一个驱动
可以模拟mysql的各种操作
相当于是让gorm来调用sqlmock来模拟测试
*/

/*
gorm在执行语句的时候会默认开一个事务
理论上让GORM执行
INSET XXX;

实际上
BEGIN;
INSERT XXX;
COMMIT;
*/
