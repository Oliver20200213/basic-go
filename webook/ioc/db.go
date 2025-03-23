package ioc

import (
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

//func InitDB() *gorm.DB {
//	//db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/webook"))
//	//连接到k8s的mysql，需要使用k8s-mysql-service.yaml中metadata的name以及port来连接
//	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:3309)/webook"))
//	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
//	if err != nil {
//		//只会再初始化过程中 panic
//		//panic相当于证个goroutine结束
//		//一段初始化过程出错，应用就不要启动了
//		panic(err)
//	}
//
//	//初始化建表语句
//	err = dao.InitTable(db)
//	if err != nil {
//		panic(err)
//	}
//
//	return db
//}

// 使用viper
func InitDB(l logger.LoggerV1) *gorm.DB {
	//dsn := viper.GetString("db.mysql.dsn")
	////viper.GetDuration("")  // 1s,1m,1h这种
	////viper.GetFloat64()   // 注意精度问题
	//fmt.Println("dsn:", dsn)

	// 注意：只有在初始化的过程中才会读取配置
	type Config struct {
		DSN string `yaml:"dsn"`

		//// 有人的做法会拆分DSN,不是很推荐
		//// localhost:13316
		//Addr string
		//// root
		//Username string
		//// root
		//Password string
		//// webook
		//DBName string

	}
	//// 设置默认值的方式2：利用结构体，常用这种方式
	//var cfg Config = Config{
	//	DSN: "root:root@tcp(localhost:3306)/webook_default",
	//}

	var cfg Config
	//err := viper.UnmarshalKey("db.mysql", &cfg) // 注意remote这里不支持key的切割，也就是不支持db.mysql，需要改下dev里面改下
	err := viper.UnmarshalKey("db", &cfg) // 注意这里不支持db.mysql需要改下dev里面改下去掉.mysql
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		// 缺了一个writer
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{ // 日志的的配置
			// 慢查询阈值，只有执行时间超过这个阈值，才会使用
			// 50ms 100ms都是比较合适的阈值
			// SQL 查询必然要求命中索引，最好就是走一次磁盘IO
			// 一次磁盘IO 是不到10ms
			SlowThreshold:             time.Millisecond * 10, // 只有超过10ms的才会输出到日志
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true, // 这里设置成true之后 sql语句中的value就会变成?
			LogLevel:                  glogger.Info,
		}),
	})
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

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}

/*
type DoSomething interface {
	DoABC() string
}

type DoSomethingFunc func() string

func (d DoSomethingFunc) DoABC() string {
	return d()
}
*/
