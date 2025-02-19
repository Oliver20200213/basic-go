package integration

import (
	"basic-go/webook/internal/web"
	"basic-go/webook/ioc"
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer() // gin.Engine，并注册了路由
	rdb := ioc.InitRedis()    // 初始化redis
	testCases := []struct {
		name string

		// 需要考虑准备数据
		before func(t *testing.T)
		// 以及验证数据数 据库的数据对不对，你redis的数据对不对
		after func(t *testing.T)

		reqBody  string
		wantCode int
		//wantBody string
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				// 不需要，也就是redis什么数据也没有
			},
			after: func(t *testing.T) {
				// 创建一个1秒超时的上下文，常用与需要控制执行时间或取消操作的场景,例如http请求 数据库查询, 并发任务等
				// 这里是如果超过1s 则释放redis数据的连接资源
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 你需要清理数据，否则连续执行两次就会失败
				val, err := rdb.GetDel(ctx, "phone_code:login:15512345678").Result()
				cancel()
				assert.NoError(t, err)
				// 你的验证码是6位
				assert.True(t, len(val) == 6)
			},
			reqBody: `{
	"phone":"15512345678"
}`,
			wantCode: 200,
			//			wantBody: `{
			//	"code":0,
			//	"msg":"发送成功",
			//}`,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				// 这个手机号码已经有一个验证码了
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 这里模拟第一次发送已经过了30s
				_, err := rdb.Set(ctx, "phone_code:login:15512345678", "123456", time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 你需要清理数据，否则连续执行两次就会失败
				// "phone_code:%s:%s"
				// GetDel获取之后删除
				val, err := rdb.GetDel(ctx, "phone_code:login:15512345678").Result()
				cancel()
				assert.NoError(t, err)
				// 你的验证码是6位,没有被覆盖，还是123456（before中是发送频繁，也就是发送失败，也就是验证码还是原来的）
				assert.Equal(t, "123456", val)
			},
			reqBody: `{
	"phone":"15512345678"
}`,
			wantCode: 200,
			//			wantBody: `{
			//	"code":0,
			//	"msg":"发送成功",
			//}`,
			wantBody: web.Result{
				Msg: "发送太频繁，请稍后再试",
			},
		},
		{
			name: "系统错误", // 模拟传入一个key没有过期时间
			before: func(t *testing.T) {
				// 这个手机号码已经有一个验证码了，但是过期时间为0 也就是没有过期时间
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 这里模拟第一次发送已经过了30s
				_, err := rdb.Set(ctx, "phone_code:login:15512345678", "123456", 0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 你需要清理数据，否则连续执行两次就会失败
				// "phone_code:%s:%s"
				val, err := rdb.GetDel(ctx, "phone_code:login:15512345678").Result()
				cancel()
				assert.NoError(t, err)
				// 你的验证码是6位,没有被覆盖，还是123456（before中是发送频繁，也就是发送失败，也就是验证码还是原来的）
				assert.Equal(t, "123456", val)
			},
			reqBody: `{
	"phone":"15512345678"
}`,
			wantCode: 200,
			//			wantBody: `{
			//	"code":0,
			//	"msg":"发送成功",
			//}`,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "手机号码为空", // 手机校验错误
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `{
	"phone":""
}`,
			wantCode: 200,
			//			wantBody: `{
			//	"code":0,
			//	"msg":"发送成功",
			//}`,
			wantBody: web.Result{
				Code: 4,
				Msg:  "输入有误",
			},
		},
		{
			name: "数据格式错误", // 需要分开来去比
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `{
	"phone":,
}`,
			wantCode: 400,
			//wantBody: web.Result{
			//	Code: 4,
			//	Msg:  "输入有误",
			//},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			// 准备request
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			// 设置json格式
			req.Header.Set("Content-Type", "application/json")
			// 准备request
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			// 需要分开比较，
			// 数据格式错误会导致json解析的时候panic
			if resp.Code != 200 {
				return
			}
			var webRes web.Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err) // 两种写法：assert.NoError(t, err)

			assert.Equal(t, tc.wantBody, webRes)

			tc.after(t)

		})
	}
}

/*
补充：
使用web.Result来进行比较而不是直接使用json字符串进行比较的原因：
1.json格式问题：
JSON 数据的格式可能会有所不同，导致它们看起来相同但实际并不完全相等。JSON 的格式化问题可能包括：
空格、换行符、缩进、顺序等
由于 JSON 的序列化规则和解析规则可能不一致，导致比较结果不一致 例如，你的 tc.wantBody 可能在某
些地方有额外的空格、换行符，或者是字段的顺序不同
2.就算是去除了不必要的空格换行符等让格式都一样，json还有编码顺序的问题
JSON 对象的字段顺序不影响 JSON 的语义，但是在进行字符串比较时，如果字段的顺序不同，字符串也会不同。
因此，即使内容相同，如果顺序不同，比较结果也会是失败的。

什么是json的编码顺序：
在 JSON 规范中，对于对象（由大括号 {} 包围的键值对集合），字段的顺序是 无关紧要的。这意味着，
对于同样的 JSON 对象，字段的排列顺序可以不同，只要键（key）和值（value）匹配，语义上它们是相同的。
例如：
{
	"name":"oliver",
	"age":35
}
{
	"age":35
	"name":"oliver",
}
这两个 JSON 对象的内容是相同的，唯一的区别在于字段顺序。尽管顺序不同，它们的语义和结构完全一致
为什么字段顺序会影响字符串比价？
当你将 JSON 对象转换为字符串时，不同的 JSON 序列化器（如 Go 的 json.Marshal）可能会以不同的
顺序输出对象的字段。这就是为什么你在测试时比较两个 JSON 字符串时，它们可能会由于顺序不同而导致比较失败。
例如：
jsonStr1 := `{"name": "Alice", "age": 25}`
jsonStr2 := `{"age": 25, "name": "Alice"}`
在直接进行字符串比较时，jsonStr1 和 jsonStr2 是不相等的，因为字段的顺序不同。但是，这两个字符串在语义上是相同的。
那该如何处理JSON的顺序问题？
常用的方式是使用结构体比较，也就是将json字符串解析成为结构体类型，然后比较结构体，而不是比较json字符串

这两种写法的区别：
err = json.NewDecoder(resp.Body).Decode(&webRes)
按需解码而不会将所有数据一次性读取到内存中，也就是只会读取一遍body
内存占用较少，适合处理大数据流
err = json.Unmarshal(resp.Body.Bytes(), &webRes)
需要先把body全部读到内存中，然后才会解析body，也就是会遍历两次body
内存占用较高，特别是大数据量的时候，可能会导致内存溢出，适合小数据量


总结：
集成测试和单元测试的不同点：
1.集成测试需要手动准备数据
2.集成测试在测试之后，是要验证数据是否符合预期的
	- 如果是使用的第三方依赖，就需要验证，比如说redis、mysql之类的
	- 如果是调用别人的接口，就不需要验证，只需要验证它返回的响应即可
	例如用到了redis，如果你在测试中向 Redis 写入了一个键值对，在测试结束后需要检测redis中是否存在这个键值对，并且值是否正确
	如果你在测试中向 MySQL 插入了一条记录，你需要查询数据库，验证这条记录是否存在且字段值是否正确。

	- 为什么要区分这两种情况
		第三方依赖（如 Redis、MySQL）：
			这些是你系统的一部分，你直接操作它们，因此你有责任确保它们的状态是正确的。
			如果这些依赖中的数据不符合预期，可能会导致你的系统行为异常。
		外部服务（别人的接口）：
			这些服务不在你的控制范围内，你无法直接访问或验证它们的内部状态。
			你只能通过它们的接口返回的响应来判断它们的行为是否符合预期。
3.注意永远不要让测试用例之间相互依赖
	例如这里发送台频繁的用例完全可以依赖于发送成功的用例（发送成功的用例执行第二次的时候就是发送频繁）
	但是不要这么做！！！互相之间有依赖关系的测试用例非常难维护！！！后续一段重构就会各种崩溃，你都改不过来！！！
*/
