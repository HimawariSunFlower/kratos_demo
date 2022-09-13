package server

import (
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"veigit-system/api/user/service/v1"
	"veigit-system/app/user/service/internal/conf"
	"veigit-system/app/user/service/internal/service"
	"veigit-system/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, ac *conf.JWT, s *service.UserService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			jwt.Server(func(token *jwt2.Token) (interface{}, error) {
				return []byte(ac.Secret), nil
			}, jwt.WithSigningMethod(jwt2.SigningMethodHS256), jwt.WithClaims(func() jwt2.Claims {
				return &auth.MyClaims{}
			})),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterUserServiceServer(srv, s)
	return srv
}
