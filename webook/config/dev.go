//go:build !k8s

//没有k8s，这个编译标签
/*
可以支持其他标签的应用：
例如：go:build dev   go:build test   go:build e2e
*/

package config

var Config = config{
	DB: DBConfig{
		//本地连接
		DSN: "localhost:13316",
	},
	Redis: RedisConfig{
		//本地连接
		Addr: "localhost:6379",
	},
}
