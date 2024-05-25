package auth

import (
	"crypto/sha256"
	"errors"
	dto "expenses/dto_expenses"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type Repository interface {
	GetUser(userName, hashPassword string) (*dto.User, error)
}

const (
	salt       = "uiwm67#hdnl4*"
	tokenTTL   = 12 * time.Hour
	signingKey = "yei#926&6%hfu*1k&j"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func hashPassword(pass string) string {
	hash := sha256.New()
	hash.Write([]byte(pass))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))

}

func generateToken(repo Repository, name, pass string) (string, error) {
	user, err := repo.GetUser(name, hashPassword(pass))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.Id,
	})
	return token.SignedString([]byte(signingKey))
}

func parseToken(inputToken string) (int, error) {
	token, err := jwt.ParseWithClaims(inputToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return -1, errors.New("invalid token claims")
	}
	return claims.UserId, nil
}
