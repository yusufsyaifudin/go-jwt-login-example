package user

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/yusufsyaifudin/go-jwt-login-example/internal/pkg/model"
	"github.com/yusufsyaifudin/go-jwt-login-example/pkg/http"
)

/**
 * @apiDefine MiddlewareAuthTokenCheck
 * @apiHeader {String} Authorization Must using Bearer access token.
 * @apiHeaderExample {json} Header-Example:
 *     {
 *       "Authorization": "Bearer your-access-token"
 *     }
 *
 * @apiParamExample {json} Request-Example:
 *     {
 *       "access_token": "your-access-token"
 *     }
 */
func (handler *HandlerConfig) MiddlewareAuthTokenCheck(next http.Handler) http.Handler {
	return func(parent context.Context, req http.Request) http.Response {
		var accessToken string

		// get access token from header
		headerAuthorization := req.RawRequest().Header.Get("Authorization")
		headerAuthorization = strings.TrimSpace(headerAuthorization)

		headerPart := strings.Split(headerAuthorization, " ")
		if len(headerPart) < 2 {
			headerPart = []string{"", ""}
		}

		if strings.ToLower(headerPart[0]) == "bearer" {
			accessToken = headerPart[1]
		}

		// if not exist on header, try using body parameter
		if accessToken == "" {
			body, err := ioutil.ReadAll(req.RawRequest().Body)
			if err != nil {
				return http.NewJsonResponse(500, map[string]interface{}{
					"error": map[string]interface{}{
						"message": fmt.Sprintf("%s: %s", "error when reading the request body", err.Error()),
					},
				})
			}

			// copy twice to make sure body can be re-binding after middleware
			body1 := ioutil.NopCloser(bytes.NewBuffer(body))
			body2 := ioutil.NopCloser(bytes.NewBuffer(body))

			var form struct {
				AccessToken string `json:"access_token" form:"access_token"`
			}

			// binding the data using body 1
			req.RawRequest().Body = body1
			req.Bind(&form)

			accessToken = form.AccessToken

			// set copied body to raw request body again
			req.RawRequest().Body = body2
		}

		jwtPayload, err := handler.Auth.ValidateToken(accessToken, handler.ServerSecretKey)
		if err != nil {
			return http.NewJsonResponse(403, map[string]interface{}{
				"error": map[string]interface{}{
					"message": fmt.Sprintf("%s: %s", "error when validating access token", err.Error()),
				},
			})
		}

		// jwtPayload.ID
		sqlGetUser := `SELECT * FROM users WHERE id = ? LIMIT 1;`

		// check user in database
		user := &model.User{}
		handler.DB.Raw(user, sqlGetUser, jwtPayload.ID)
		if user == nil || user.ID == 0 {
			return http.NewJsonResponse(401, map[string]interface{}{
				"error": map[string]interface{}{
					"message": "cannot continue this request since user is not found with this token",
				},
			})
		}

		req.SetUser(user)

		// run the wrapped handler
		return next(parent, req)
	}
}
