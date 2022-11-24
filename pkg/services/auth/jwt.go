package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtClaims struct {
	jwt.StandardClaims
	UserId    int64     `json:"userId"`
	Email     string    `json:"email"`
	LastLogin time.Time `json:"lastLogin"`
}

func (c *JwtClaims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}
	if c.UserId == 0 {
		return errors.New("userId is empty")
	}
	if c.Email == "" {
		return errors.New("email is empty")
	}
	if c.LastLogin.IsZero() {
		return errors.New("lastLogin is zero")
	}
	return nil
}

func GenerateJwt(jwtKey string, id int64, email string) (string, error) {
	claims := &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
		UserId:    id,
		Email:     email,
		LastLogin: time.Now(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

func VerifyJwt(jwtKey string, tokenString string) (*JwtClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error parsing jwt")
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, errors.New("error parsing jwt")
	}
	if err := claims.Valid(); err != nil {
		return nil, err
	}
	return claims, nil
}
