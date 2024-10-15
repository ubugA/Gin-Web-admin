package jwtoken

import (
	"time"

	"gin-api-admin/internal/proposal"

	"github.com/golang-jwt/jwt/v5"
)

var _ Token = (*token)(nil)

type Token interface {
	i()
	Sign(jwtInfo proposal.SessionUserInfo, expireDuration time.Duration) (tokenString string, err error)
	Parse(tokenString string) (*claims, error)
}

type token struct {
	secret string
}

type claims struct {
	proposal.SessionUserInfo
	jwt.RegisteredClaims
}

func New(secret string) Token {
	return &token{
		secret: secret,
	}
}

func (t *token) i() {}

func (t *token) Sign(sessionUserInfo proposal.SessionUserInfo, expireDuration time.Duration) (tokenString string, err error) {
	claims := claims{
		sessionUserInfo,
		jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
	}

	tokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(t.secret))
	return
}

func (t *token) Parse(tokenString string) (*claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
