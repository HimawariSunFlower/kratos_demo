package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"veigit-system/api/user/service/v1"
	"veigit-system/app/user/service/internal/conf"
	"veigit-system/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           uint
	Phone        int
	Username     string
	RealName     string
	PasswordHash string `copier:"Password"`
}

type UserLogin struct {
	Phone    int
	Username string
	Token    string
}

type UserUpdate struct {
	Phone    int
	Username string
	Password string
}

func hashPassword(pwd string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func verifyPassword(hashed, input string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input)); err != nil {
		return false
	}
	return true
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByPhone(ctx context.Context, phone int) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByID(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUserById(ctx context.Context, id uint) error
}

type UserUsecase struct {
	ur   UserRepo
	jwtc *conf.JWT
	log  *log.Helper
}

type Profile struct {
	ID        uint
	Username  string
	Bio       string
	Image     string
	Following bool
}

func NewUserUsecase(ur UserRepo, logger log.Logger, jwtc *conf.JWT) *UserUsecase {
	return &UserUsecase{ur: ur, log: log.NewHelper(logger), jwtc: jwtc}
}

func (uc *UserUsecase) generateToken(userID uint) string {
	ret, _ := auth.GenerateToken(userID, uc.jwtc.Secret)
	return ret
}

func (uc *UserUsecase) Register(ctx context.Context, in *User) (*UserLogin, error) {
	u := &User{
		Phone:        in.Phone,
		Username:     in.Username,
		RealName:     in.RealName,
		PasswordHash: hashPassword(in.PasswordHash),
	}
	if err := uc.ur.CreateUser(ctx, u); err != nil {
		return nil, err
	}
	return &UserLogin{
		Phone:    in.Phone,
		Username: in.Username,
		Token:    uc.generateToken(u.Id),
	}, nil
}

func (uc *UserUsecase) Login(ctx context.Context, phone int, password string) (*UserLogin, error) {
	if phone <= 0 {
		return nil, v1.ErrorParamError("phone cannot be empty")
	}
	u, err := uc.ur.GetUserByPhone(ctx, phone)
	if err != nil {
		return nil, v1.ErrorSqlError(err.Error())
	}
	if !verifyPassword(u.PasswordHash, password) {
		return nil, errors.Unauthorized("user", "login failed")
	}

	return &UserLogin{
		Phone:    u.Phone,
		Username: u.Username,
		Token:    uc.generateToken(u.Id),
	}, nil
}

func (uc *UserUsecase) GetCurrentUser(ctx context.Context) (*User, error) {
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		return nil, v1.ErrorContentMissing("jwt Claims error")
	}
	my := claims.(*auth.MyClaims)
	u, err := uc.ur.GetUserByID(ctx, my.Uid)
	if err != nil {
		return nil, v1.ErrorSqlError(err.Error())
	}
	return u, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, uu *UserUpdate) (*UserLogin, error) {
	u, err := uc.ur.GetUserByUsername(ctx, uu.Username)
	if err != nil {
		return nil, v1.ErrorSqlError(err.Error())
	}
	u.Phone = uu.Phone
	u.PasswordHash = hashPassword(uu.Password)
	u, err = uc.ur.UpdateUser(ctx, u)
	if err != nil {
		return nil, v1.ErrorSqlError(err.Error())
	}
	return &UserLogin{
		Phone:    u.Phone,
		Username: u.Username,
		Token:    uc.generateToken(u.Id),
	}, nil
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id uint) (*User, error) {
	u, err := uc.ur.GetUserByID(ctx, id)
	if err != nil {
		return nil, v1.ErrorSqlError(err.Error())
	}
	return u, nil
}

func (uc *UserUsecase) GetUserByUsername(ctx context.Context, name string) (*User, error) {
	u, err := uc.ur.GetUserByUsername(ctx, name)
	if err != nil {
		return nil, v1.ErrorSqlError(err.Error())
	}
	return u, nil
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id uint) error {
	err := uc.ur.DeleteUserById(ctx, id)
	if err != nil {
		return v1.ErrorSqlError(err.Error())
	}
	return nil
}
