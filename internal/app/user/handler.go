package user

import (
	"github.com/yusufsyaifudin/go-jwt-login-example/pkg/auth"
	"github.com/yusufsyaifudin/go-jwt-login-example/pkg/db"
)

type HandlerConfig struct {
	ServerSecretKey string
	DB              db.Query
	Auth            auth.Auth
}

func NewUserHandler(serverSecretKey string, db db.Query, auth auth.Auth) *HandlerConfig {
	return &HandlerConfig{
		ServerSecretKey: serverSecretKey,
		DB:              db,
		Auth:            auth,
	}
}
