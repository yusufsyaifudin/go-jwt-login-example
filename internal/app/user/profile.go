package user

import (
	"context"

	"github.com/yusufsyaifudin/go-jwt-login-example/pkg/http"
)

/**
 * @api {get} /user/profile Profile
 * @apiVersion 1.0.0
 * @apiName Get Profile
 * @apiGroup User
 *
 * @apiDescription Get user profile, based on authentication header.
 *
 * @apiHeader {String} Authorization Authorization value, using format `Bearer {user-jwt-access-token}.
 */
func (handler *HandlerConfig) ProfileUserHandler(ctx context.Context, req http.Request) http.Response {
	user := req.User()

	return http.NewJsonResponse(200, map[string]interface{}{
		"user": map[string]interface{}{
			"id":            user.ID,
			"name":          user.Name,
			"username":      user.Username,
			"registered_at": user.CreatedAt.Unix(),
		},
	})
}
