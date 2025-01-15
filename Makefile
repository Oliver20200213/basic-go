.PHONY:mock
mock:
	@mockgen --source=E:/gowork/src/basic-go/webook/internal/service/user.go -package=svcmocks -destination=E:/gowork/src/basic-go/webook/internal/service/mocks/user.mock.go
	@mockgen --source=E:/gowork/src/basic-go/webook/internal/service/code.go -package=svcmocks -destination=E:/gowork/src/basic-go/webook/internal/service/mocks/code.mock.go
	@ go mod tidy
