package ioc

import (
	"basic-go/webook/config"
	"basic-go/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	//db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/webook"))
	//连接到k8s的mysql，需要使用k8s-mysql-service.yaml中metadata的name以及port来连接
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:3309)/webook"))
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		//只会再初始化过程中 panic
		//panic相当于证个goroutine结束
		//一段初始化过程出错，应用就不要启动了
		panic(err)
	}

	//初始化建表语句
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}
