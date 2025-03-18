package ioc

//// InitUserService 单独配置logger，userservice中独有的logger
//func InitUserService(repo repository.UserRepository) service.UserService {
//	l, err := zap.NewDevelopment() // 这里直接是创建了一个新的Logger
//	if err != nil {
//		panic(err)
//	}
//	return service.NewUserService(repo, l)
//	// 直接传入了Logger，在service中可以直接使用这个Logger
//	// 而不是调用zap.L()取全局的Logger，所以不需要zap.ReplaceGlobals()
//}
