package service

import (
	"context"
	"veigit-system/api/system/interface/v1"
	"veigit-system/app/system/interface/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedInterfaceServer

	ac *biz.AuthUseCase
	uc *biz.UserUseCase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.UserUseCase) *GreeterService {
	return &GreeterService{uc: uc}
}

func (s *GreeterService) Login(ctx context.Context, req *v1.LoginReq) (reply *v1.LoginReply, err error) {
	return s.ac.Login(ctx, req)
}

func (s *GreeterService) Logout(ctx context.Context, req *v1.LogoutReq) (reply *v1.LogoutReply, err error) {
	//todo jwt logout
	return &v1.LogoutReply{}, nil
}
func (s *GreeterService) Register(ctx context.Context, req *v1.RegisterReq) (reply *v1.RegisterReply, err error) {
	return s.ac.Register(ctx, req)
}
