package http

import (
	"net/http"

	"github.com/yusufsyaifudin/go-jwt-login-example/internal/pkg/model"
)

type Request interface {
	ContentType() string
	Bind(out interface{}) error
	GetParam(key string) string
	RawRequest() *http.Request
	User() *model.User // get the current user
	SetUser(user *model.User)
}
