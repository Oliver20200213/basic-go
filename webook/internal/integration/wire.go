//go:build wireinject

package integration

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	"basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,

		//初始化DAO
		dao.NewUserDAO,

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewCodeRepository,
		repository.NewUserRepository,

		service.NewUserService,
		service.NewCodeService,

		// 基于内存的实现
		//memory.NewService,
		ioc.InitSMSService,

		web.NewUserHandler,

		// 中间件怎么办
		// 注册路由怎么办
		// 这个地方没有用到前面的任何东西
		//gin.Default,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)

	// 返回的是可以随便返回的
	return new(gin.Engine)
}
