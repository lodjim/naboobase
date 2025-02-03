package utils

import (
	"fmt"
	"naboobase/configs"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var jwtKey = []byte(configs.GetSecretKey())

type JwtToken struct {
	Token string `validate:"required,jwt"`
}

type Claims struct {
	Email       string `json:"email"`
	Id          string `json:"id"`
	IsSuperUser bool   `json:"is_super_user"`
	IsVerified  bool   `json:"is_verified"`
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	Email string `json:"email"`
	Id    string `json:"id"`
	jwt.StandardClaims
}

var validate = validator.New()

func CreateToken(email string, id string, isVerified bool, isSuperUser bool) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email:       email,
		Id:          id,
		IsSuperUser: isSuperUser,
		IsVerified:  isVerified,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CreateRefreshToken(email string, id string) (string, error) {
	duration := configs.GetExpirationTime()
	expirationTime := time.Now().Add(time.Duration(duration) * time.Minute)
	claims := &RefreshTokenClaims{
		Email: email,
		Id:    id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetClaims(c *gin.Context) (*Claims, error) {
	claims := &Claims{}
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return claims, nil
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func VerifyJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func VerifyRefreshJWT(tokenStr string) (*RefreshTokenClaims, error) {
	claims := &RefreshTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
