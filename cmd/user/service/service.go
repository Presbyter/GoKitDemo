package service

import (
	repo "GoKitDemo/repositories"
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
)

type Service interface {
	Login(ctx context.Context, loginType int32, value string, deviceType int32, deviceCode string) (*repo.User, error)
}

func New(logger log.Logger /*, ints, chars metrics.Counter*/) Service {
	var svc Service
	{
		svc = NewUserService()
		svc = LoggingMiddleware(logger)(svc)
		//svc = InstrumentingMiddleware(ints, chars)(svc)
	}
	return svc
}

type userService struct {
	userRepo *repo.UserRepo
}

func NewUserService() Service {
	return &userService{
		userRepo: repo.NewUserRepo(),
	}
}

func (u *userService) Login(ctx context.Context, loginType int32, value string, devictType int32, deviceCode string) (*repo.User, error) {
	//switch loginType {
	//case 0: // Phone
	//u.userRepo.GetByPhone(strings.TrimSpace(value))
	//case 1: // WeChat
	//case 2: // QQ
	//case 3: // Apple
	//default:
	//	return nil, errors.New("not sport LoginType")
	//}

	//panic("not implement")
	return &repo.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		PhoneNum:     "13333333333",
		QQOpenID:     "",
		WeChatOpenID: "",
		AppleID:      "",
		NickName:     "",
		Gender:       0,
		Avatar:       "",
		Status:       0,
	}, nil
}
