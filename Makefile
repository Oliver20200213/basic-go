.PHONY:mock
mock:
	@mockgen --source=E:/gowork/src/basic-go/webook/internal/service/user.go -package=svcmocks -destination=E:/gowork/src/basic-go/webook/internal/service/mocks/user.mock.go
	@mockgen --source=E:/gowork/src/basic-go/webook/internal/service/code.go -package=svcmocks -destination=E:/gowork/src/basic-go/webook/internal/service/mocks/code.mock.go
	@mockgen --source=E:/gowork/src/basic-go/webook/internal/repository/user.go -package=repomocks -destination=E:/gowork/src/basic-go/webook/internal/repository/mocks/user.mock.go
	@mockgen --source=E:/gowork/src/basic-go/webook/internal/repository/code.go -package=repomocks -destination=E:/gowork/src/basic-go/webook/internal/repository/mocks/code.mock.go
	@mockgen --source=E:/gowork/src/basic-go/webook/internal/repository/dao/user.go -package=daomocks -destination=E:/gowork/src/basic-go/webook/internal/repository/dao/mocks/user.mock.go
	@mockgen --source=E:/gowork/src/basic-go/webook/internal/repository/cache/user.go -package=cachemocks -destination=E:/gowork/src/basic-go/webook/internal/repository/cache/mocks/user.mock.go
	@mockgen -package=redismocks -destination=E:/gowork/src/basic-go/webook/internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
