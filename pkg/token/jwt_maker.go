package token

import (
	"time"

	"github.com/AGODOVALOV/grader/pkg/token/config"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(cfg *config.Config) (Maker, error) {
	duration, err := time.ParseDuration(cfg.Duration)
	if err != nil {
		return nil, err
	}
	return &JWTMaker{
		secretKey:     cfg.JWTSecret,
		tokenDuration: duration}, nil
}

// CreateToken creates a new JWT token
func (J *JWTMaker) CreateToken(userID int64, username string) (string, *Payload, error) {
	payload, err := NewPayload(userID, username, J.tokenDuration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(J.secretKey))
	return token, payload, err
}

// VerifyToken verifies a JWT token
func (J *JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(J.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, payload.Valid()
}
