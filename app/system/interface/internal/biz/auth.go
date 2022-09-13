package biz

import (
	"context"
	"errors"
	v1 "veigit-system/api/system/interface/v1"
	"veigit-system/app/system/interface/internal/conf"
	"veigit-system/pkg/middleware/auth"
)

var (
	ErrLoginFailed = errors.New("login failed")
)

type AuthUseCase struct {
	key      string
	userRepo UserRepo
}

func NewAuthUseCase(conf *conf.Auth, userRepo UserRepo) *AuthUseCase {
	return &AuthUseCase{
		key:      conf.ApiKey,
		userRepo: userRepo,
	}
}

func (receiver *AuthUseCase) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginReply, error) {

	// get user
	user, err := receiver.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, v1.ErrorLoginFailed("user not found: %s", err.Error())
	}
	// check permission(password blacklist etc...)
	err = receiver.userRepo.VerifyPassword(ctx, user, req.Password)
	if err != nil {
		return nil, v1.ErrorLoginFailed("password not match")
	}

	signedString, err := auth.GenerateToken(uint(user.Id), receiver.key)

	if err != nil {
		return nil, v1.ErrorLoginFailed("generate token failed: %s", err.Error())
	}
	return &v1.LoginReply{
		Token: signedString,
	}, nil
}

func (receiver *AuthUseCase) Register(ctx context.Context, req *v1.RegisterReq) (*v1.RegisterReply, error) {

	// check username
	_, err := receiver.userRepo.FindByUsername(ctx, req.Username)
	if !errors.Is(err, ErrUserNotFound) {
		return nil, v1.ErrorRegisterFailed("username already exists")
	}
	// create user
	user, err := NewUser(req.Username, req.Password)
	if err != nil {
		return nil, v1.ErrorRegisterFailed("create user failed: %s", err.Error())
	}
	// save user
	err = receiver.userRepo.Save(ctx, &user)
	if err != nil {
		return nil, v1.ErrorRegisterFailed("save user failed: %s", err.Error())
	}
	return &v1.RegisterReply{
		Id: uint64(user.Id),
	}, nil
}
