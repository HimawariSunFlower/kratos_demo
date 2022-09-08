package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	userv1 "veigit-system/api/user/service/v1"
	"veigit-system/app/system/interface/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewUserServiceClient)

// Data .
type Data struct {
	log *log.Helper
	uc  userv1.UserServiceClient
}

// todo db
// NewData .
func NewData(c *conf.Data, uc userv1.UserServiceClient, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{uc: uc}, cleanup, nil
}

func NewUserServiceClient(ac *conf.Auth) userv1.UserServiceClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("192.168.20.10:9000"),
		//grpc.WithDiscovery(r),
		//grpc.WithMiddleware(
		//	tracing.Client(tracing.WithTracerProvider(tp)),
		//	recovery.Recovery(),
		//	jwt.Client(func(token *jwt2.Token) (interface{}, error) {
		//		return []byte(ac.ServiceKey), nil
		//	}, jwt.WithSigningMethod(jwt2.SigningMethodHS256)),
		//),
	)
	if err != nil {
		panic(err)
	}
	c := userv1.NewUserServiceClient(conn)
	return c
}
