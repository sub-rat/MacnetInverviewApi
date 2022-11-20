package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type authCustomClaims struct {
	Id uint `json:"id"`
	jwt.StandardClaims
}

// GenerateJwtToken generates the token which has validity of 48 hrs
// and SIGNING_KEY provided from environment
func GenerateJwtToken(id uint) string {
	claims := &authCustomClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	fmt.Printf("%v %v", tokenString, err)
	return tokenString
}

// VerifyJwtToken Parse and verify the Token and return the Token and error
func VerifyJwtToken(r *gin.Context) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SIGNING_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ExtractToken returns the Token string  from Authorization header
func ExtractToken(r *gin.Context) string {
	bearToken := r.GetHeader("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
