package service

import (
	"context"
	"veigit-system/api/user/service/v1"
	"veigit-system/app/user/service/internal/biz"
	"veigit-system/pkg/tools"
)

func (s *UserService) GetUserList(ctx context.Context, req *v1.GetUserListReq) (reply *v1.UserReply, err error) {
	//todo 是否需要？
	return &v1.UserReply{}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserReq) (reply *v1.UserReply, err error) {
	u, err := s.uc.GetUserByID(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}
	ret := &v1.UserReply{User: &v1.User{}}
	tools.Copy(u, ret.User)
	return ret, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *v1.UpdateUserReq) (reply *v1.UserReply, err error) {
	input := &biz.UserUpdate{}
	tools.Copy(req.Data, input)
	u, err := s.uc.UpdateUser(ctx, input)
	if err != nil {
		return nil, err
	}
	ret := &v1.UserReply{User: &v1.User{}}
	tools.Copy(u, ret.User)
	return ret, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *v1.DeleteUserReq) (reply *v1.UserReply, err error) {
	err = s.uc.DeleteUser(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{}, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, req *v1.GetUserByUsernameReq) (reply *v1.UserReply, err error) {
	u, err := s.uc.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	ret := &v1.UserReply{User: &v1.User{}}
	tools.Copy(u, ret.User)
	return ret, nil
}
