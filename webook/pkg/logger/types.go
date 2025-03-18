package logger

// 风格一：这种风格，要求用户必须在 msg 里面留好占位符
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func LoggerExample() {
	var l Logger
	phone := "155xxxx5678"
	l.Info("用户未注册，手机号码是%s,", phone)
}

// 风格二： zap就是这个风格，认为日志里面的参数都是有名字的
type LoggerV1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
}

type Field struct {
	Key   string
	Value any
}

func LoggerV1Example() {
	var l Logger
	phone := "155xxxx5678"
	l.Info("用户未注册", Field{
		Key:   "phone",
		Value: phone,
	})
}

type LoggerV2 interface {
	// args必须是偶数，并且按照key-value来组织
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func LoggerV2Example() {
	var l Logger
	phone := "155xxxx5678"
	l.Info("用户未注册", "phone", phone)
}
