package zapx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 实现phone的脱敏

type MyCore struct {
	zapcore.Core
}

func (c MyCore) Write(entry zapcore.Entry, fds []zapcore.Field) error {
	for _, fd := range fds {
		if fd.Key == "phone" {
			phone := fd.String
			fd.String = phone[:3] + "****" + phone[7:]
		}
	}
	return c.Core.Write(entry, fds)
}

func MaskPhone(key string, Phone string) zap.Field {
	Phone = Phone[:3] + "****" + Phone[7:]
	return zap.Field{
		Key:    key,
		String: Phone,
	}
}
