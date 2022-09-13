package data

import (
	"context"
	"veigit-system/app/user/service/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

type User struct {
	gorm.Model
	Email        string `gorm:"size:500"`
	Username     string `gorm:"size:500"`
	RealName     string `gorm:"size:500"`
	Phone        int
	PasswordHash string `gorm:"size:500"`
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) CreateUser(ctx context.Context, u *biz.User) error {
	user := User{
		Phone:        u.Phone,
		Username:     u.Username,
		RealName:     u.RealName,
		PasswordHash: u.PasswordHash,
	}
	rv := r.data.db.Create(&user)
	return rv.Error
}

func (r *userRepo) UpdateUser(ctx context.Context, in *biz.User) (rv *biz.User, err error) {
	u := new(User)
	err = r.data.db.Where("username = ?", in.Username).First(u).Error
	if err != nil {
		return nil, err
	}
	err = r.data.db.Model(&u).Updates(&User{
		Phone:        in.Phone,
		PasswordHash: in.PasswordHash,
	}).Error
	return &biz.User{
		Id:           u.ID,
		Phone:        u.Phone,
		Username:     u.Username,
		RealName:     u.RealName,
		PasswordHash: u.PasswordHash,
	}, nil
}

func (r *userRepo) GetUserByPhone(ctx context.Context, phone int) (rv *biz.User, err error) {
	u := new(User)
	result := r.data.db.Where("phone = ?", phone).First(u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("user", "not found by email")
	}
	if result.Error != nil {
		return nil, err
	}
	return &biz.User{
		Id:           u.ID,
		Phone:        u.Phone,
		Username:     u.Username,
		RealName:     u.RealName,
		PasswordHash: u.PasswordHash,
	}, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id uint) (rv *biz.User, err error) {
	u := new(User)
	err = r.data.db.Where("id = ?", id).First(u).Error
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Id:           u.ID,
		Phone:        u.Phone,
		Username:     u.Username,
		RealName:     u.RealName,
		PasswordHash: u.PasswordHash,
	}, nil
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (rv *biz.User, err error) {
	u := new(User)
	err = r.data.db.Where("username = ?", username).First(u).Error
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Id:           u.ID,
		Phone:        u.Phone,
		Username:     u.Username,
		RealName:     u.RealName,
		PasswordHash: u.PasswordHash,
	}, nil
}

func (r *userRepo) DeleteUserById(ctx context.Context, id uint) (err error) {
	u := new(User)
	err = r.data.db.Where("id = ?", id).Delete(u).Error
	return err
}
