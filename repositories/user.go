package repositories

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	PhoneNum     string         `json:"phone_num"`
	QQOpenID     string         `json:"qq_open_id"`
	WeChatOpenID string         `json:"wechat_open_id"`
	AppleID      string         `json:"apple_id"`
	NickName     string         `json:"nick_name"`
	Gender       GenderEnum     `json:"gender"`
	Avatar       string         `json:"avatar"`
	Status       UserStatusEnum `json:"status"`
}

type GenderEnum uint8

const (
	GenderEnum_Man GenderEnum = iota
	GenderEnum_Woman
)

type UserStatusEnum uint8

const (
	UserStatusEnum_Active UserStatusEnum = iota
	UserStatusEnum_Disable
)

//---------------------------------------------------------------------------------------------------------------------

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo() *UserRepo {
	return &UserRepo{db: mysqlDB}
}

func (u *UserRepo) Get(id uint) (*User, error) {
	entity := User{}
	err := u.db.First(&entity, id).Error
	return &entity, err
}

func (u *UserRepo) GetByPhone(phone string) (*User, error) {
	entity := User{}
	err := u.db.First(&entity, User{PhoneNum: phone}).Error
	return &entity, err
}

func (u *UserRepo) GetByQQ(openID string) (*User, error) {
	entity := User{}
	err := u.db.First(&entity, User{QQOpenID: openID}).Error
	return &entity, err
}

func (u *UserRepo) GetByWeChat(openID string) (*User, error) {
	entity := User{}
	err := u.db.First(&entity, User{WeChatOpenID: openID}).Error
	return &entity, err
}

func (u *UserRepo) Create(entity User) error {
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()

	return u.db.Create(&entity).Error
}
