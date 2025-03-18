package logger

import "sync"

/*
不使用依赖注入，使用全局变量
*/

var gl LoggerV1
var lMutex sync.RWMutex

func SetGlobalLogger(l LoggerV1) {
	lMutex.Lock()
	defer lMutex.Unlock()
	gl = l
}

func L() LoggerV1 {
	lMutex.RLock()
	g := gl
	lMutex.RUnlock()
	return g
}

var GL LoggerV1 = &NopLogger{}
