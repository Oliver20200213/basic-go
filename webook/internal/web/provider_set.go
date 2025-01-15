package web

// UserProviderSet 在这里定义如何构造我的UserHandle
// 我的这个UserHandler里面只专注于提供HTTP服务
// 你愿意怎么初始化就怎么初始化
// 某个功能，以及如何用要分离
// 某个类型，以及某个类型的实例，怎么构造，也可以分离
// 放在这里是不合适的， 万一我不用wire呢
// 不太赞成的用法
//var UserProviderSet = wire.NewSet(
//	ioc.InitDB, ioc.InitRedis,
//
//	//初始化DAO
//	dao.NewUserDAO,
//
//	cache.NewUserCache,
//	cache.NewCodeCache,
//
//	repository.NewCodeRepository,
//	repository.NewUserRepository,
//
//	service.NewUserService,
//	service.NewCodeService,
//
//	ioc.InitSMSService,
//
//	NewUserHandler,
//)
