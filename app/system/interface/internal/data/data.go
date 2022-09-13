package data

import (
	"context"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	userv1 "veigit-system/api/user/service/v1"
	"veigit-system/app/system/interface/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewDB, NewAD, NewUserServiceClient)

// Data .
type Data struct {
	log *log.Helper
	uc  userv1.UserServiceClient
	db  *gorm.DB
	ad  *gormadapter.Adapter
}

// NewData . todo db
func NewData(c *conf.Data, uc userv1.UserServiceClient, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{uc: uc}, cleanup, nil
}

//	func NewUserServiceClient(ac *conf.Auth) userv1.UserServiceClient {
//		conn, err := grpc.DialInsecure(
//			context.Background(),
//			grpc.WithEndpoint("192.168.20.10:9000"), //todo 服务发现
//			//grpc.WithDiscovery(r),
//			grpc.WithMiddleware(
//				//tracing.Client(tracing.WithTracerProvider(tp)),
//				recovery.Recovery(),
//				jwt.Client(func(token *jwt2.Token) (interface{}, error) {
//					return []byte(ac.ServiceKey), nil
//				}, jwt.WithSigningMethod(jwt2.SigningMethodHS256)),
//			),
//		)
//		if err != nil {
//			panic(err)
//		}
//		c := userv1.NewUserServiceClient(conn)
//		return c
//	}
func NewUserServiceClient(ac *conf.Auth, r registry.Discovery) userv1.UserServiceClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///veigit.user.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
			jwt.Client(func(token *jwt2.Token) (interface{}, error) {
				return []byte(ac.ServiceKey), nil
			}, jwt.WithSigningMethod(jwt2.SigningMethodHS256)),
		),
	)
	if err != nil {
		panic(err)
	}
	c := userv1.NewUserServiceClient(conn)
	return c
}

func NewDB(c *conf.Data) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err) //,"failed to connect database")
	}
	return db
}

func NewAD(db *gorm.DB) *gormadapter.Adapter {
	a, _ := gormadapter.NewAdapterByDBWithCustomTable(db, &gormadapter.CasbinRule{}, "test_casbin_rule")
	return a
}
