package jwt

import (
	"context"
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/constant"
)

type JwtInterface interface {
	GenerateToken(id int64, userName string) (token string, err error)
	ParseToken(tokenString string) (jwtClaims *JwtClaims, err error)
	GetUser(ctx context.Context) (*JwtClaims, error)
}

type GenerateToken struct {
	secret string
	expire time.Duration
	issuer string
}

func NewGenerateToken() (*GenerateToken, error) {
	var (
		secret string
		expire time.Duration
		issuer string
		err    error
	)
	if secret, err = conf.GetJwtSecret(); err != nil {
		return nil, err
	}

	issuer = conf.GetJwtIssuer()

	if expire, err = conf.GetJwtExpirationTime(); err != nil {
		return nil, err
	}
	return &GenerateToken{
		secret: secret,
		expire: expire,
		issuer: issuer,
	}, nil
}

type JwtClaims struct {
	UserID   int64  `json:"userId"`
	UserName string `json:"userName"`
	*jwtv5.RegisteredClaims
}

func newJwtClaims(userID int64, userName, issuer string, expire time.Duration) *JwtClaims {
	now := time.Now()
	return &JwtClaims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: &jwtv5.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwtv5.NewNumericDate(now.Add(expire)),
			IssuedAt:  jwtv5.NewNumericDate(now),
			NotBefore: jwtv5.NewNumericDate(now),
		},
	}
}

func (j *GenerateToken) GenerateToken(id int64, userName string) (token string, err error) {
	jwtClaims := newJwtClaims(id, userName, j.issuer, j.expire)

	claims := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, jwtClaims)

	token, err = claims.SignedString([]byte(j.secret))
	if err != nil {
		return "", errors.New("generate token failed")
	}
	return token, nil
}

// ParseToken 解析token
func (j *GenerateToken) ParseToken(tokenString string) (jwtClaims *JwtClaims, err error) {
	jwtClaims = &JwtClaims{}
	token, err := jwtv5.ParseWithClaims(tokenString, jwtClaims, func(token *jwtv5.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, errors.New("parse token failed")
	}
	if claims, ok := token.Claims.(*JwtClaims); ok {
		return claims, nil
	}
	return nil, errors.New("parse token failed")
}

func (j *GenerateToken) GetUser(ctx context.Context) (*JwtClaims, error) {
	cl := ctx.Value(constant.UserContextKey)
	if cl == nil {
		return nil, errors.New("get jwt claims by ctx failed")
	}
	jwtClaims, ok := cl.(*JwtClaims)
	if !ok {
		return nil, errors.New("get jwt claims by ctx failed")
	}
	return jwtClaims, nil
}
