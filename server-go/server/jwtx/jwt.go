package jwtx

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var key string = "server-go"

// 生成token
func GenToken(param map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	for k, v := range param {
		claims[k] = v
	}
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // 7天有效期
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 验证token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("valid failed")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token.Claims failed")
	}
	return claims, nil
}
