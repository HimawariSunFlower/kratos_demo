package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"math/rand"
	v1 "veigit-system/api/system/interface/v1"
	"veigit-system/pkg/tools"
)

var (
	ErrPasswordInvalid = errors.New("password invalid")
	ErrUsernameInvalid = errors.New("username invalid")
	ErrUserNotFound    = errors.New("user not found")
)

type User struct {
	Id       uint64
	Username string
	Password string
}

func NewUser(username string, password string) (User, error) {
	// check username
	if len(username) <= 0 {
		return User{}, ErrUsernameInvalid
	}
	// check password
	if len(password) <= 0 {
		return User{}, ErrPasswordInvalid
	}

	return User{
		Id:       rand.Uint64(), //todo uuid
		Username: username,
		Password: password,
	}, nil
}

type UserRepo interface {
	Find(ctx context.Context, id uint64) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	Save(ctx context.Context, u *User) error

	VerifyPassword(ctx context.Context, u *User, password string) error
}

type UserUseCase struct {
	repo   UserRepo
	log    *log.Helper
	authUc *AuthUseCase
}

func NewUserUseCase(repo UserRepo, logger log.Logger, authUc *AuthUseCase) *UserUseCase {
	log := log.NewHelper(log.With(logger, "module", "usecase/interface"))
	return &UserUseCase{
		repo:   repo,
		log:    log,
		authUc: authUc,
	}
}

func (uc *UserUseCase) Logout(ctx context.Context, u *User) error {
	return nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, uid uint64) (*v1.GetUserReply, error) {
	user, err := uc.repo.Find(ctx, uid)
	if err != nil {
		return nil, v1.ErrorUnknownError(err.Error())
	}
	ret := &v1.GetUserReply{}
	tools.Copy(user, ret)
	return ret, nil
}
