package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yusufsyaifudin/go-jwt-login-example/internal/pkg/model"
	"github.com/yusufsyaifudin/go-jwt-login-example/pkg/auth"
	"github.com/yusufsyaifudin/go-jwt-login-example/pkg/http"
)

/**
 * @api {post} /user/login Login
 * @apiVersion 1.0.0
 * @apiName Login
 * @apiGroup User
 *
 * @apiDescription User login
 *
 * @apiParam (Request body) {String} username Username of registered user
 * @apiParam (Request body) {String} password User password
 */
func (handler *HandlerConfig) LoginUserHandler(ctx context.Context, req http.Request) http.Response {
	form := &struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}{}

	if err := req.Bind(form); err != nil {
		return http.NewJsonResponse(500, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("fail when binding the payload: %s", err.Error()),
			},
		})
	}

	if strings.TrimSpace(form.Username) == "" {
		return http.NewJsonResponse(400, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "username cannot be empty",
			},
		})
	}

	if strings.TrimSpace(form.Password) == "" {
		return http.NewJsonResponse(400, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "password cannot be empty",
			},
		})
	}

	// check user in database
	user := &model.User{}
	handler.DB.Raw(user, "SELECT * FROM users WHERE username = ? LIMIT 1", form.Username)
	if user == nil || user.ID == 0 {
		return http.NewJsonResponse(404, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "user not found",
			},
		})
	}

	if !CheckPasswordHash(form.Password, user.Password) {
		return http.NewJsonResponse(401, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "wrong password",
			},
		})
	}

	// if found, then check hashing password
	tokenPayload := &auth.Payload{
		ID:        fmt.Sprintf("%d", user.ID),
		Username:  user.Username,
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		ExpiredAt: time.Now().Add(5 * time.Hour).Unix(),
	}

	accessToken, err := handler.Auth.GenerateToken(tokenPayload, handler.ServerSecretKey)
	if err != nil {
		return http.NewJsonResponse(422, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("fail generating access token: %s", err.Error()),
			},
		})
	}

	return http.NewJsonResponse(200, map[string]interface{}{
		"access_token": accessToken,
		"user": map[string]interface{}{
			"id":            user.ID,
			"name":          user.Name,
			"username":      user.Username,
			"registered_at": user.CreatedAt.Unix(),
		},
	})
}
