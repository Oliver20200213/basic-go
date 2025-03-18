package logger

import "go.uber.org/zap"

/*
使用风格二来封装zap
缺陷是这里有一个参数的转换，会引入不必要的内存分配和CPU消耗
这种也叫做适配器模式


什么是适配器模式
就是将一个类型适配到另外一个接口（将zap的Logger适配到我们自己的接口）
注意：这里说的额是类型，也就是说A适配到B，那么A可以是一个接口，也可以是一个具体类型
典型的使用场景：
- 你有两个用途类似，但是细节上有所不同的接口，然后将一个接口（类型）适配到另外一个接口
- 你因为版本升级，出现两个不同的接口，那么就得把老的几口适配到新的接口，也可以把新的接口适配到老的接口
  这取决于你怎用

注意：装饰器模式，一直都是同一个接口，而适配器模式，必然是不同的接口
*/

type ZapLogger struct {
	l *zap.Logger
}

func NewZapLogger(l *zap.Logger) LoggerV1 {
	return &ZapLogger{l: l}
}

func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.l.Debug(msg, z.toZapFields(args)...)
}

func (z *ZapLogger) Info(msg string, args ...Field) {
	z.l.Info(msg, z.toZapFields(args)...)
}

func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.l.Warn(msg, z.toZapFields(args)...)
}

func (z *ZapLogger) Error(msg string, args ...Field) {
	z.l.Error(msg, z.toZapFields(args)...)
}

// 将自定义的Field类型切片转换成zap.Field的类型切片
func (z *ZapLogger) toZapFields(args []Field) []zap.Field {
	res := make([]zap.Field, 0, len(args)) // 生成一个长度为0 容量为len(args)的切片
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
		//zap.Any()是将任意类型的值转为zap.Field,第一个参数为key，第二个参数为value
	}
	return res
}
