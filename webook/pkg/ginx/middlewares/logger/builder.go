package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type MiddlewareLoggerBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewMiddlewareLoggerBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareLoggerBuilder {
	return &MiddlewareLoggerBuilder{
		loggerFunc: fn,
	}
}

func (b *MiddlewareLoggerBuilder) AllowReqBody() *MiddlewareLoggerBuilder {
	b.allowReqBody = true
	return b
}

func (b *MiddlewareLoggerBuilder) AllowRespBody() *MiddlewareLoggerBuilder {
	b.allowRespBody = true
	return b
}

func (b *MiddlewareLoggerBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 { // 这个1024可以做成参数
			url = url[:1024]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			// url也可能很长，可以不打全部，
			Url: url,
		}
		if b.allowReqBody && ctx.Request.Body != nil {
			// 直接忽略 error 不影响程序运行
			//body, _ := io.ReadAll(ctx.Request.Body)
			body, _ := ctx.GetRawData() // 本质是调用io.ReadAll(ctx.Request.Body)来读取数据
			// 注意：要再把body放回去，body是一个io.ReadCloser，可以想象成一个流，这个body读完就没了
			// io这个包里面基本上就是涉及流的操作，读完就没了
			// 放回去就会有性能损耗
			reader := io.NopCloser(bytes.NewReader(body)) // 这里用bytes.NewBuffer()也可以用，只要是实现ReadCloser就行
			ctx.Request.Body = reader
			// 另一种写法，不太常用，可能中间件中有用过
			//ctx.Request.GetBody = func()(io.ReadCloser, error){
			//	return reader, nil
			//}
			//if len(body) > 1024 {
			//	body = body[:1024]
			//}
			// 这其实是一个很消耗cpu和内存的操作
			// 因为会引起复制
			al.ReqBody = string(body)
		}

		if b.allowRespBody {
			// 将ctx中原来的gin.ResponseWriter，替换成自己的responseWriter
			ctx.Writer = responseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}

		// 防止意外panic而导致日志没有打印
		defer func() {
			al.Duration = time.Since(start).String()
			//al.Duration = time.Now().Sub(start)

			b.loggerFunc(ctx, al)
		}()

		// 执行到业务逻辑
		ctx.Next()

		//b.loggerFunc(ctx, al)
	}
}

/*
在 gin 框架中，response.Body 并不能像 request.Body 那样直接获取，因为 gin.Context.Writer
只会将数据写入 HTTP 响应流，而不会存储响应内容。
通过装饰的方式劫持，gin.ResponseWriter，从而拦截Write()方法，记录响应的Body
重写gin.ResponseWriter的WriterHeader、Write、WriteString这三个方法
*/

// 通过装饰器的方式获取request.Body的数据
type responseWriter struct {
	al                 *AccessLog
	gin.ResponseWriter // 使用组合的方式， 如果是只需要装饰部分方法，那么就用组合，如果是需要装饰所有的就用非组合的方式
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// 为什么会有Write和WriteString看起来是重复了
// Write是处理二进制数据的，如JSON 图片等
// WriteString是处理纯文本字符串的，如HTML 文本响应等
// 所以需要同时拦截这两个方法
func (w responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriteString(data string) (int, error) {
	// 要不要控制长度呢？ 按需配置
	w.al.RespBody = data
	return w.ResponseWriter.WriteString(data)
}

type AccessLog struct {
	// HTTP 请求的方法
	Method string
	// URL 整个请求
	Url string
	// 可以根据自己的需求增加

	Duration string

	// ReqBody,RespBody 比较复杂，比较多 要考虑是否要打印
	ReqBody  string
	RespBody string
	Status   int
}
