//go:build wireinject

package main

import (
	"basic-go/practice/internal/repository"
	"basic-go/practice/internal/repository/cache"
	"basic-go/practice/internal/repository/dao"
	"basic-go/practice/internal/service"
	"basic-go/practice/internal/web"
	"basic-go/practice/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,

		dao.NewUserDao,

		cache.NewCodeCache,
		cache.NewUserCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,

		ioc.InitSms,
		web.NewUserHandler,

		ioc.InitMiddleware,
		ioc.InitWebServer,
	)
	return new(gin.Engine)
}
