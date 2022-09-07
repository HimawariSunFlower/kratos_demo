package data

import (
	"context"
	"veigit-system/pkg/tools"

	"github.com/go-kratos/kratos/v2/log"
	usV1 "veigit-system/api/user/service/v1"
	"veigit-system/app/system/interface/internal/biz"
)

var _ biz.UserRepo = (*userRepo)(nil)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (rp *userRepo) Find(ctx context.Context, id int64) (*biz.User, error) {
	//TODO implement me
	panic("implement me")
}

func (rp *userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	user, err := rp.data.uc.GetUserByUsername(ctx, &usV1.GetUserByUsernameReq{UserName: username})
	if err != nil {
		return nil, biz.ErrUserNotFound
	}
	return &biz.User{
		Id:       user.User.Id,
		Username: user.User.UserName,
	}, nil
}

func (rp *userRepo) Save(ctx context.Context, u *biz.User) error {
	toVal := &usV1.User{}
	tools.Copy(u, toVal)
	_, err := rp.data.uc.UpdateUser(ctx, &usV1.UpdateUserReq{
		Data: toVal,
	})
	return err
}

func (rp *userRepo) VerifyPassword(ctx context.Context, u *biz.User, password string) error {
	return nil
	//TODO implement me
	panic("implement me")
}
