package auth

import (
	"kratos-realworld/internal/pkg/tools"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestGenerateToken(t *testing.T) {
	tk := GenerateToken("secret", 11)
	spew.Dump(tk)
	panic("tk")
}

type User struct {
	Id    int64
	Phone int64
}
type User1 struct {
	Id    int64
	Phone int64
}

func TestName(t *testing.T) {
	u := &User1{Id: 1, Phone: 123}
	to := &User{}
	tools.Copy(u, to)
}
