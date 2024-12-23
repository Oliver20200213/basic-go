//go:build k8s

// 使用k8s 这个编译标签
package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-mysql:3309)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:6379",
	},
}
