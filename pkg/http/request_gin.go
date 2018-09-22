package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yusufsyaifudin/go-jwt-login-example/internal/pkg/model"
)

// WrapGin wraps a Handler and turns it into gin compatible handler
// This method should be called with a fresh ctx
func WrapGin(parent context.Context, handler Handler) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// create span
		ctx := context.Background()
		defer ctx.Done()

		// create request and run the handler
		var req = newGinRequest(ginContext)
		resp := handler(ctx, req)

		if resp == nil {
			ginContext.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "internal_server_error",
					"message": "nil response",
					"data":    nil,
				},
			})
			return
		}

		// get the body first
		body, err := resp.Body()
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "internal_server_error",
					"message": err.Error(),
					"data":    nil,
				},
			})
			return
		}
		// then write header
		for k, v := range resp.Header() {
			for _, h := range v {
				ginContext.Writer.Header().Add(k, h)
			}
		}

		ginContext.Writer.Header().Add("Content-Type", resp.ContentType())
		ginContext.Writer.WriteHeader(resp.StatusCode())

		// the last is writing the body
		ginContext.Writer.Write(body)
	}
}

type ginRequest struct {
	context *gin.Context
	user    *model.User
}

func newGinRequest(context *gin.Context) (request Request) {
	request = &ginRequest{
		context: context,
	}
	return
}

func (ginRequest *ginRequest) Bind(out interface{}) error {
	return ginRequest.context.Bind(out)
}

func (ginRequest *ginRequest) GetParam(key string) string {
	return ginRequest.context.Param(key)
}

func (ginRequest *ginRequest) Header() http.Header {
	return ginRequest.context.Request.Header
}

func (ginRequest *ginRequest) ContentType() string {
	return ginRequest.context.ContentType()
}

func (ginRequest *ginRequest) RawRequest() *http.Request {
	return ginRequest.context.Request
}

func (ginRequest *ginRequest) User() *model.User {
	return ginRequest.user
}

func (ginRequest *ginRequest) SetUser(user *model.User) {
	ginRequest.user = user
}
