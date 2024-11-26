package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"context"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// u domain.User 可以不用指针，如果使用指针需要进行判空 if u ==nil{}
	//很有可能分配到栈上去  没有内存逃逸

	//需要考虑加密放到哪里
	//然后就是存起来
	return svc.repo.Create(ctx, u)
}
