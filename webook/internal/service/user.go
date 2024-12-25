package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账户/邮箱或密码错误")
var ErrInvalidUserID = errors.New("无效的用户ID")

type UserService struct {
	repo *repository.UserRepository
	//redis *redis.Client //错误实践 不是很严谨的方式
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
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	//然后就是存起来
	return svc.repo.Create(ctx, u)

}

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	//先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFind) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		//DEBUG日志
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) Edit(ctx context.Context, u domain.User) error {
	u, err := svc.repo.FindById(ctx, u.Id)
	if errors.Is(err, repository.ErrUserNotFind) {
		return ErrInvalidUserID
	}
	if err != nil {
		return err
	}

	return nil
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	//错误实践
	////第一个念头是
	//val, err := svc.redis.Get(ctx, fmt.Sprintf("user:info:%d", id)).Result()
	//if err != nil {
	//	return domain.User{}, err
	//}
	//var u domain.User
	//err = json.Unmarshal([]byte(val), &u)
	//if err != nil {
	//	return u, err
	//}
	//接下来，就是从数据库中查找

	//最佳实践
	u, err := svc.repo.FindById(ctx, id)
	return u, err

}
