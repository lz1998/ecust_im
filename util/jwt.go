package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lz1998/ecust_im/config"
	"github.com/lz1998/ecust_im/model/user"
)

func GenerateJwtTokenString(ecustUser *user.EcustUser) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["expire"] = time.Now().Add(time.Hour * time.Duration(24)).Unix()
	claims["userId"] = ecustUser.UserId
	claims["nickname"] = ecustUser.Nickname
	claims["password"] = ecustUser.Password
	claims["status"] = ecustUser.Status
	claims["createdAt"] = ecustUser.CreatedAt.Unix()
	claims["updatedAt"] = ecustUser.UpdatedAt.Unix()
	token.Claims = claims
	return token.SignedString(config.JwtSecret)
}

func JwtParseUser(tokenString string) (*user.EcustUser, error) {
	if tokenString == "" {
		return nil, errors.New("no token is found in Authorization Bearer")
	}
	claims := make(jwt.MapClaims)
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	ecustUser := &user.EcustUser{
		UserId:    int64(claims["userId"].(float64)),
		Nickname:  claims["nickname"].(string),
		Password:  claims["password"].(string),
		Status:    int64(claims["updatedAt"].(float64)),
		UpdatedAt: time.Unix(int64(claims["updatedAt"].(float64)), 0),
		CreatedAt: time.Unix(int64(claims["createdAt"].(float64)), 0),
	}
	return ecustUser, nil
}
