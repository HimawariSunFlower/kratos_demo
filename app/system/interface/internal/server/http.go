package server

import (
	"context"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwt2 "github.com/golang-jwt/jwt/v4"
	v1 "veigit-system/api/system/interface/v1"
	"veigit-system/app/system/interface/internal/conf"
	"veigit-system/app/system/interface/internal/service"
	"veigit-system/pkg/casbin"
	"veigit-system/pkg/middleware/auth"
)

func NewWhiteListMatcher() selector.MatchFunc {

	whiteList := make(map[string]struct{})
	whiteList["/system.interface.v1.Interface/Login"] = struct{}{}
	whiteList["/system.interface.v1.Interface/Register"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, ac *conf.Auth, casbin2 *conf.Casbin, ad *gormadapter.Adapter, greeter *service.GreeterService, logger log.Logger) *http.Server {
	m, _ := model.NewModelFromFile(casbin2.Model)
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			selector.Server(
				jwt.Server(func(token *jwt2.Token) (interface{}, error) {
					return []byte(ac.ApiKey), nil
				}, jwt.WithSigningMethod(jwt2.SigningMethodHS256), jwt.WithClaims(&auth.MyClaims{})),
				casbin.Server(casbin.WithCasbinModel(m), casbin.WithCasbinPolicy(ad)),
			).Match(NewWhiteListMatcher()).Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterInterfaceHTTPServer(srv, greeter)
	return srv
}
