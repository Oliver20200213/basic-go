package main

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	//db := initDB()
	//server := initServer()
	//rdb := initRedis()
	//u := initUser(db, rdb)
	//u.RegisterRoutes(server)

	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "这是hello Go页面！")
	//})
	initViperRemote()

	//initViperWatch()
	keys := viper.AllKeys()
	fmt.Println("keys:", keys)
	setting := viper.AllSettings()
	fmt.Println("setting:", setting)

	server := InitWebServer()

	//server.RunTLS()  https
	server.Run(":8080")

}

// viper使用
// 方式1：
func initViper() {
	// 配置文件的名字，但是不包含文件扩展名
	// 不包含.go, .yaml之类的后缀
	viper.SetConfigName("dev")
	// 告诉viper 我的配置用的是yaml格式
	// 实际使用中，有很多格式，JSON, XML, YAML, TOML, INI
	viper.SetConfigType("yaml")
	// 指定配置文件的存储路径，当前工作下的config子目录
	viper.AddConfigPath("./config")
	//// 可以有多个，会依次扫描
	//viper.AddConfigPath("/ect/webook")
	// 读取配置到viper里面，或者可以理解为加载到内存里面
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	//// 可以创建多套配置实例
	//otherViper := viper.New()
	//otherViper.SetConfigName("myjson")
	//otherViper.SetConfigType("json")
	//otherViper.AddConfigPath("./config")
}

// viper没有分辨key是否存在的方法，也既是没法区别这个key到底有没有，
// 如果没有只会返回对应字段的零值 这里dsn是string，返回的就是""空字符串
func initViperV1() {
	// 设置默认值方式1：
	//viper.SetDefault("db.mysql.dsn", "root@root@tcp(localhost:3306)/mysql")

	// 直接指定文件路径
	//viper.SetConfigFile("config/dev.yaml") // 是从go的working directory开始定位的
	// 使用绝对路径也可以，不推荐 因为如果项目换电脑就会失效
	//viper.SetConfigFile("E:\\gowork\\src\\basic-go\\webook\\config\\dev.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperReader() {
	viper.SetConfigType("yaml")
	cfg := `
db.mysql:
  dsn: "root:root@tcp(localhost:13316)/webook"


redis:
  addr: "localhost:6379"`
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg))) // 将字符串转换成io.Reader使用bytes.NewReader()
	// 可以从文件里面读，网络中读，内存中读
	if err != nil {
		panic(err)
	}
}

// 区分不同的测试环境
// 使用ide直接调用会报错，
// 需要配置下ide里面的Program arguments选项： --config=config/dev.yaml
// 命令行中使用方式： go run . --config=config/dev.yaml
func initViperV2() {
	/*
		config: 参数的名称key
		config/config.yaml: 参数config对应的默认参数值value
		指定配置文件路径：参数config的描述信息（帮助信息）
	*/
	cfile := pflag.String("config",
		"config/config.yaml",
		"指定配置文件路径")
	pflag.Parse()               // 解析，要先解析在使用
	viper.SetConfigFile(*cfile) // 注意解引用
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperRemote() {
	/*
			需要到github中下载etcdctl
			windows：找个release版本，直接下载，然后放到go环境变量workspace的bin下面即可
			linux：
				git clone https://github.com/etcd-io/etcd.git
				cd etcd
				git tag 找到要安装的relaese版本
				git checkout -b intall v3.5.9 切入到这个版本
				cd etcdctl
		        go install .  直接安装（同样会安装到go的workspace目录下bin）
			执行下etcdctl 查看是否安装成功

			将配置文件存储到etcd中
			linux/unix的命令：
			（windows下不可用，如果使用gitbash，会有问题会将C:/Program Files/Git/webook当做key而不是/webook）
			cd config
			etcdctl --endpoints=127.0.0.1:12379 put /webook "$(<dev.yaml)"
			put 存储命令  /webook 键key  $(<dev.yaml) 值value  Bash语法等价于cat dev.yaml
			向etcd中存储一个键值对
			windows:(powershell)
			etcdctl --endpoints=127.0.0.1:12379 put /webook "$(Get-Content dev.yaml -Raw)"
			查看存储的信息：
			etcdctl --endpoints=127.0.0.1:12379 get /webook
			获取数据执行返回的结构：
			PS E:\gowork\src\basic-go\webook\config> etcdctl --endpoints=127.0.0.1:12379 get /webook
			/webook
			db:
			  dsn: root:root@tcp(localhost:13316)/webook


			redis:
			  addr: localhost:6379
	*/
	viper.SetConfigType("yaml")
	err := viper.AddRemoteProvider("etcd3",
		"127.0.0.1:12379", // etcd的连接地址
		// 存储到etcd中的key，可以理解为通过webook和其他使用etcd的区别出来
		"/webook")
	if err != nil {
		panic(err)
	}
	// 远程配置中心没法监听变动在变动的时候更新，只能每一次都实时的读取数据
	//err = viper.WatchRemoteConfig()
	//if err != nil {
	//	panic(err)
	//}
	//// 只适用于文件的变化，远程不适用
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	fmt.Println(in.Name, in.Op)
	//})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func initViperWatch() {
	cfile := pflag.String("config",
		"config/dev.yaml",
		"指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	// 实时监听配置变更
	viper.WatchConfig() // 注意顺序，放在read前面即可
	// 只能告诉你文件变了，不能告诉你哪里变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 比较好的设计，他会在in里面告诉你变更前的数据，和变更后的数据
		// 更好的设计是，他会告诉你差异
		fmt.Println(in.Name, in.Op)
		// 需要重新读才能拿到新的配置
		fmt.Println(viper.GetString("db.dsn"))
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

//func initDB() *gorm.DB {
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

//func initServer() *gin.Engine {
//	server := gin.Default()
//
//	server.Use(func(c *gin.Context) { //先注册先执行后注册后执行
//		fmt.Println("这是第一个middleware")
//	})
//	server.Use(func(ctx *gin.Context) {
//		fmt.Println("这是第二个middleware")
//	})
//
//	//引入redis限流
//	redisClient := redis.NewClient(&redis.Options{
//		//对应k8s-redis-service中metadata的name以及port
//		Addr: config.Config.Redis.Addr,
//	})
//	//1分钟中内只能允许100个请求
//	server.Use(ratelimit.NewBuilder(redisClient, time.Minute, 100).Build())
//
//	server.Use(cors.New(cors.Config{ //use方法注册middleware，use作用于全部路由
//		//AllowOrigins: []string{"http://127.0.0.1:3000"}, //允许的访问的源
//		//AllowOrigins: []string{"*"},  允许所有请求，不建议这种方式，之前前端是可以的现在前端在严格模式下不生效
//		//AllowMethods: []string{"POST","GET"}, //如果不写，默认则是都支持
//		AllowHeaders:     []string{"Content-Type", "Authorization"},
//		ExposeHeaders:    []string{"x-jwt-token"}, //不加这个前端是拿不到x-jwt-token的，意思是我给你的你才能拿到
//		AllowCredentials: true,                    //是否允许你带cookie之类的东西
//
//		//可以根据origin进行动态判断(注销掉上面的AllowOrigins配置，二选一配置)
//		AllowOriginFunc: func(origin string) bool {
//			//如果origin包含http://127.0.0.1就允许访问
//			if strings.Contains(origin, "http://127.0.0.1") {
//				//你的开发环境
//				return true
//			}
//			return strings.Contains(origin, "公司的域名")
//		},
//		MaxAge: 12 * time.Hour, //profile的有效期
//	}))
//
//	//session实现步骤1
//	//单实例可以使用memstore,将数据存储到内存中(常用于测试或开发环境中)
//	//多实例使用redistore
//	//第一个参数是authentication key（加密key） 第二个是encryption key（加密value）
//	//store := memstore.NewStore([]byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"),
//	//	[]byte("pueKIkHTQsCIMa1N7mmzkTN6NmmHjIOP"))
//	//server.Use(sessions.Sessions("ssid", store))
//
//	//store, err := redis.NewStore(16, //第一个参数表示最大空闲连接数量，实际中随便写，16,32都行
//	//	"tcp", "localhost:6379", "", //协议 登录地址和端口 密码
//	//	[]byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"), //authentication key 指身份认证  32位或64位的key都可以
//	//	[]byte("pueKIkHTQsCIMa1N7mmzkTN6NmmHjIOP")) //Encryption 是指数据加密
//	//if err != nil {
//	//	panic(err)
//	//}
//	//server.Use(sessions.Sessions("ssid", store))
//
//	//根据面向接口编程 实现自定义sqlxstore
//	//mystore := &sqlx_store.Store{}
//	//server.Use(sessions.Sessions("ssid", mystore))
//
//	//store := cookie.NewStore([]byte("secret"))   //存储的地方，数据存储到cookie
//	//server.Use(sessions.Sessions("ssid", store)) //mysession是cookie中的名字，store是值
//	//session实现步骤3
//	//server.Use(middlewares.NewLoginMiddlewareBuilder().Build())
//	//链式调用,session最好的实现
//	//server.Use(middlewares.NewLoginMiddlewareBuilder().
//	//	IgnorePaths("/users/login").
//	//	IgnorePaths("/users/signup.lua").Build())
//
//	//版本1
//	////忽略sss路径
//	//middlewares.IgnorePaths = []string{"sss"}
//	//server.Use(middlewares.CheckLogin())
//	////又有一个server不能忽略sss这个路径,此时v1版本无法实现
//	//server1 := gin.Default()
//	//server1.Use(middlewares.CheckLogin())
//
//	//使用JWT
//	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
//		IgnorePaths("/users/signup.lua").
//		IgnorePaths("/users/login").
//		IgnorePaths("/users/login_sms/code/send").
//		IgnorePaths("/users/login_sms").
//		Build())
//
//	return server
//}

//func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
//	ud := dao.NewUserDAO(db)
//	uc := cache.NewUserCache(rdb) //初始化用户缓存 user cache
//	repo := repository.NewUserRepository(ud, uc)
//	svc := service.NewUserService(repo)
//	codeCache := cache.NewCodeCache(rdb)
//	codeRepo := repository.NewCodeRepository(codeCache)
//	smsSvc := memory.NewService()
//	codeSvc := service.NewCodeService(codeRepo, smsSvc)
//	u := web.NewUserHandler(svc, codeSvc)
//	return u
//}

//func initRedis() redis.Cmdable {
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: config.Config.Redis.Addr,
//	})
//	return redisClient
//}
