package auth

// Payload is a data carried by JWT token
type Payload struct {
	ID        string `json:"id"`       // required, id of this user
	Username  string `json:"username"` // required, name of this user
	IssuedAt  int64  `json:"iss"`      // token creation date, epoch time in seconds value (10 character)
	NotBefore int64  `json:"nbf"`      // token valid start date, if token used before this time, it will contain error, epoch time in seconds value (10 character)
	ExpiredAt int64  `json:"exp"`      // token expiration date, epoch time in seconds value (10 character)
}

// Auth is an higher abstraction level of authorization method.
// ValidateToken method: to check if a token is valid or not, and
// GenerateToken method: to generate token based on jwt payload
// By this interface, you can easily change the JWT 3rd party library if it doesn't meet your needs.
type Auth interface {
	GenerateToken(payload *Payload, secretKey string) (token string, err error)
	ValidateToken(token string, secretKey string) (payload *Payload, err error)
}
