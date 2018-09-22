package auth

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtPayload extends the base auth Payload, so all Payload property can be read here
type JwtPayload struct {
	Payload
}

// Jwt will implements Auth interface using library github.com/dgrijalva/jwt-go.
type Jwt struct {
}

// NewJwtAuth is like a class implementing interface Auth
func NewJwtAuth() (auth Auth) {
	auth = &Jwt{}
	return
}

// Valid is a method required by jwt.Claims set (github.com/dgrijalva/jwt-go).
// Inside this function, we will do any validation checking.
func (jwtPayload JwtPayload) Valid() error {
	if strings.TrimSpace(jwtPayload.ID) == "" {
		return fmt.Errorf("id must contains value")
	}

	if strings.TrimSpace(jwtPayload.Username) == "" {
		return fmt.Errorf("name must contains value")
	}

	if len(fmt.Sprintf("%d", jwtPayload.IssuedAt)) != 10 {
		return fmt.Errorf("iat must in epoch time contained 10 character length")
	}

	if jwtPayload.IssuedAt > time.Now().Unix() {
		return fmt.Errorf("token used before issued")
	}

	if len(fmt.Sprintf("%d", jwtPayload.NotBefore)) != 10 {
		return fmt.Errorf("nbf must in epoch time contained 10 character length")
	}

	if jwtPayload.NotBefore < time.Now().Unix() {
		return fmt.Errorf("token is not valid yet")
	}

	if len(fmt.Sprintf("%d", jwtPayload.ExpiredAt)) != 10 {
		return fmt.Errorf("exp must in epoch time contained 10 character length")
	}

	if jwtPayload.ExpiredAt <= time.Now().Unix() {
		return fmt.Errorf("token is expired")
	}

	return nil
}

// GenerateToken will generate jwt token using inputted payload
func (authJwt *Jwt) GenerateToken(payload *Payload, secretKey string) (token string, err error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	// generate token using HS256
	jwtToken.Header = map[string]interface{}{
		"alg": jwt.SigningMethodHS256.Name,
	}

	jwtToken.Claims = JwtPayload{
		Payload: *payload,
	}

	token, err = jwtToken.SignedString([]byte(secretKey)) // sign with secret key
	return
}

// ValidateToken implements validating jwt token using secret key and return payload
func (authJwt *Jwt) ValidateToken(token string, secretKey string) (payload *Payload, err error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Don't forget to validate the alg is what you expect. We use HMAC algorithm.
	// _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC)
	// if !ok {
	// 	err = fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"])
	// 	return nil, err
	// }

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		err = fmt.Errorf("token is not valid payload")
		return
	}

	if claims.Valid() != nil {
		return nil, claims.Valid()
	}

	// convert type jwt.MapClaims to type Payload
	mapClaimBytes, err := json.Marshal(claims)
	if err != nil {
		err = fmt.Errorf("payload cannot be marshalled")
		return nil, err
	}

	payload = &Payload{}
	err = json.Unmarshal(mapClaimBytes, payload)
	if err != nil {
		err = fmt.Errorf("payload cannot be unmarshalled")
		return nil, err
	}

	// build payload based on jwt payload
	return payload, nil
}
