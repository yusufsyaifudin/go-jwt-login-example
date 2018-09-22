package auth_test

import (
	"testing"
	"time"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/smartystreets/goconvey/convey"
	"github.com/yusufsyaifudin/go-jwt-login-example/pkg/auth"
)

func TestGenerateAndValidateTokenAtBestCondition(t *testing.T) {
	t.Parallel()

	secretKey := "abc"
	authJwt := auth.NewJwtAuth()

	convey.Convey("Generate and validate token", t, func() {

		convey.Convey("When all value is good", func() {
			inputPayload := &auth.Payload{
				ID:        "1",
				Username:  "John Doe",
				IssuedAt:  time.Now().Unix(),
				NotBefore: time.Now().Unix(),
				ExpiredAt: time.Now().Add(2 * time.Minute).Unix(),
			}

			jwtToken, err := authJwt.GenerateToken(inputPayload, secretKey)
			convey.So(err, convey.ShouldBeNil)

			outputPayload, err := authJwt.ValidateToken(jwtToken, secretKey)
			convey.So(err, convey.ShouldBeNil)
			convey.So(outputPayload, convey.ShouldResemble, inputPayload)
		})

		convey.Convey("When secret key is different", func() {
			inputPayload := &auth.Payload{
				ID:        "1",
				Username:  "John Doe",
				IssuedAt:  time.Now().Unix(),
				NotBefore: time.Now().Unix(),
				ExpiredAt: time.Now().Add(2 * time.Minute).Unix(),
			}

			jwtToken, err := authJwt.GenerateToken(inputPayload, secretKey)
			convey.So(err, convey.ShouldBeNil)

			outputPayload, err := authJwt.ValidateToken(jwtToken, "cba")
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(err.Error(), convey.ShouldResemble, "signature is invalid")
			convey.So(outputPayload, convey.ShouldNotResemble, inputPayload)
		})

	})
}

func TestNewJwtAuth(t *testing.T) {
	t.Parallel()

	convey.Convey("Test initiating new jwt auth", t, func() {
		convey.Convey("Should return JWT struct and implements Auth interface", func() {
			var jwtAuth auth.Auth
			jwtAuth = auth.NewJwtAuth()
			convey.So(jwtAuth, convey.ShouldNotBeNil)
			convey.So(jwtAuth, convey.ShouldEqual, &auth.Jwt{})
		})
	})
}

func TestJwt_ValidateToken(t *testing.T) {
	t.Parallel()

	convey.Convey("Validate token", t, func() {
		secretKey := "abc"
		authJwt := auth.NewJwtAuth()

		// This occurred if GenerateToken function, at some point becouse wrong logic, returns wrong signing method
		convey.Convey("When signing method is different", func() {
			// https://stackoverflow.com/a/51472209/5489910
			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			if err != nil {
				log.Fatal(err)
			}

			claims := &jwt.StandardClaims{
				ExpiresAt: 15000,
				Issuer:    "test",
			}

			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

			tokenString, err := token.SignedString(key)
			convey.So(err, convey.ShouldBeNil)
			convey.So(token, convey.ShouldNotBeNil)

			// test the function
			jwtPayload, err := authJwt.ValidateToken(tokenString, secretKey)
			convey.So(jwtPayload, convey.ShouldBeNil)
			convey.So(err.Error(), convey.ShouldResemble, "key is of invalid type")
		})

		convey.Convey("Token is expired", func() {
			tokenString := "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjEiLCJuYW1lIjoiSm9obiBEb2UiLCJpc3MiOjE1MzY0OTA0MDgsIm5iZiI6MTUzNjQ5MDQwOCwiZXhwIjoxNTM2NDkwNTI4fQ.Qz8gVmKS6v75S8TLcyteT0H3J5_6EO6R0f6h9OmhBJ0"
			jwtPayload, err := authJwt.ValidateToken(tokenString, secretKey)
			convey.So(jwtPayload, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})

	})
}
