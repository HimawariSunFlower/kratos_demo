package service

import (
	"github.com/google/wire"
	"veigit-system/api/user/service/v1"

	"veigit-system/app/user/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserService)

type UserService struct {
	v1.UnimplementedUserServiceServer

	uc  *biz.UserUsecase
	log *log.Helper
}

func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{uc: uc, log: log.NewHelper(logger)}
}
