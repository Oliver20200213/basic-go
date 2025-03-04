package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("账户/邮箱或密码错误")

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WeChatInfo) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
}

type userService struct {
	repo repository.UserRepository
	//redis *redis.Client //错误实践 不是很严谨的方式
}

// 我用的人，只管用，怎初始化我不管， 我一点都不关心如何初始化
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
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

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	//先找用户
	u, err := svc.repo.FindByEmail(ctx, email)

	if errors.Is(err, repository.ErrUserNotFound) {
		// 用户未找到
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		// 系统错误
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

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 先查找
	// 为啥要先查找？
	// 如果没有查找直接通过insert成功说明是新用户，失败说明是已注册过的用户。理论上是可以的
	// 但是如果是1w个用户里面有9500个是已注册过的，那么就需要1w次的insert
	// 如果是先查找那么insert只需要500次，9500次只需要查询即可
	// 这个叫做快路径
	u, err := svc.repo.FindByPhone(ctx, phone)
	//判断，有没有这个用户
	if err != repository.ErrUserNotFound {
		// 绝大部分请求会来这里
		// nil 会进来这里
		// 不为 ErrUserNotFound也会进来这里
		return u, err
	}

	// 在系统资源不足，触发降级之后，不执行慢路径
	//if ctx.Value("降级") == "true" {
	//	return domain.User{}, errors.New("系统降级了")
	//}

	// 这个叫做慢路径
	// 你明确知道 没有这个用户
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	// 注册有问题但是又不是手机号码冲突，说明一定是系统错误
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 注意这里会遇到主从延迟的问题
	return svc.repo.FindByPhone(ctx, phone)

}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, info domain.WeChatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, info.OpenId)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	u = domain.User{
		WeChatInfo: info,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}

	return svc.repo.FindByWechat(ctx, info.OpenId)
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	//错误实践，在service中调用cache
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

	// 最佳实践
	// 在系统内部，基本上都是用ID的
	// 有些系统比较复杂，会用一个GUID(global unique ID)
	return svc.repo.FindById(ctx, id)

}

// 降级操作
//func PathsDownGrade(ctx context.Context, quick, slow func()) {
//	quick()
//	if ctx.Value("降级") == "true" {
//		return
//	}
//	slow()
//}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}
