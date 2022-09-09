package casbin

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"veigit-system/pkg/constant"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	jwtV4 "github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cast"
)

type contextKey string

const (
	reason           string = "FORBIDDEN"
	defaultRBACModel        = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*") || r.sub == "admin"
`
)

var (
	ErrEnforcerMissing     = errors.Forbidden(reason, "Enforcer is missing")
	ErrSecurityParseFailed = errors.Forbidden(reason, "Security Info fault")
	ErrContext             = errors.Forbidden(reason, "Context Info fault")
	ErrUnauthorized        = errors.Forbidden(reason, "Unauthorized Access")
)

type Option func(*options)

type options struct {
	model    model.Model
	policy   persist.Adapter
	enforcer *casbin.SyncedEnforcer
}

func WithCasbinModel(model model.Model) Option {
	return func(o *options) {
		o.model = model
	}
}

func WithCasbinPolicy(policy persist.Adapter) Option {
	return func(o *options) {
		o.policy = policy
	}
}

// loadRbacModel 加载RBAC模型
func loadRbacModel() (model.Model, error) {
	return model.NewModelFromString(defaultRBACModel)
}

func Server(opts ...Option) middleware.Middleware {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	if o.model == nil {
		o.model, _ = loadRbacModel()
	}

	if o.enforcer == nil {
		o.enforcer, _ = casbin.NewSyncedEnforcer(o.model, o.policy)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				allowed bool
				err     error
			)

			if o.enforcer == nil {
				return nil, ErrEnforcerMissing
			}

			claims, ok := jwt.FromContext(ctx)
			if !ok {
				return nil, ErrSecurityParseFailed
			}
			mc, ok := claims.(jwtV4.MapClaims)
			if !ok {
				return nil, ErrSecurityParseFailed
			}
			userid, ok := mc[constant.JWT_USERID]
			if !ok {
				return nil, ErrSecurityParseFailed
			}

			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, ErrContext
			}
			htr, ok := tr.(*http.Transport)
			if !ok {
				return nil, ErrContext
			}

			u := cast.ToString(userid)
			allowed, err = o.enforcer.Enforce(u, htr.Request().URL.Path, htr.Request().Method)

			if err != nil {
				return nil, errors.Forbidden(reason, err.Error())
			}
			if !allowed {
				return nil, ErrUnauthorized
			}
			return handler(ctx, req)
		}
	}
}
