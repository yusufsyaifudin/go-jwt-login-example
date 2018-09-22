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
 * @api {post} /user/register Register
 * @apiVersion 1.0.0
 * @apiName Register
 * @apiGroup User
 *
 * @apiDescription User register. This also return authentication token for the first time.
 *
 * @apiParam (Request body) {String} name Name of this user
 * @apiParam (Request body) {String} username Username of the user. This should be unique.
 * @apiParam (Request body) {String} password User password
 */
func (handler *HandlerConfig) RegisterUserHandler(ctx context.Context, req http.Request) http.Response {
	form := &struct {
		Name     string `json:"name" form:"name"`
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

	if strings.TrimSpace(form.Name) == "" {
		return http.NewJsonResponse(400, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "name cannot be empty",
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

	// check if user already exist
	user := &model.User{}
	handler.DB.Raw(user, "SELECT * FROM users WHERE username = ? LIMIT 1", form.Username)
	if user != nil && user.ID != 0 {
		return http.NewJsonResponse(400, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "user with this username already registered",
			},
		})
	}

	passwordHash, err := HashPassword(form.Password)
	if err != nil {
		return http.NewJsonResponse(422, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("fail when hashing password: %s", err.Error()),
			},
		})
	}

	var sqlInsertUser = `
		INSERT INTO users (name, username, password) VALUES (?, ?, ?) ON CONFLICT(username) DO UPDATE SET updated_at = now() RETURNING *;
	`

	// insert to db user in database
	err = handler.DB.Raw(user, sqlInsertUser, form.Name, form.Username, passwordHash)
	if err != nil {
		return http.NewJsonResponse(422, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("fail inserting user into db: %s", err.Error()),
			},
		})
	}

	// Check password hash is different or not with body json data, if different, it may because attacking.
	// If still the same, it may because race condition in request (2 or more request at one time)
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
