package server

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	salt       = "uiwm67#hdnl4*"
	tokenTTL   = 12 * time.Hour
	signingKey = "yei#926&6%hfu*1k&j"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

// validateIdExpense validate id is not empty and convert id's type string to int
func validateIdExpense(id string) (int, error) {
	var expenseID int
	var err error
	if id != "" {
		expenseID, err = strconv.Atoi(id)
		if err != nil {
			return -1, errors.New("incorrect id")
		}
	} else {
		return -1, errors.New("incorrect id")
	}
	return expenseID, nil
}

// checkExtension checks type file
func checkExtension(fileName string) (typeFile string, err error) {
	src := strings.Split(fileName, ".")
	//check name of file has extension and create a new file with random characters appended to the name
	if len(src) == 2 {
		extension := strings.ToLower(src[1])
		switch extension {
		case "doc", "pdf", "txt":
			typeFile = "document"
		case "jpg", "jpeg", "png", "gif", "raw", "svg", "bmp", "ico", "tiff", "webp":
			typeFile = "image"
		case "mp4", "webm", "mov", "avi", "flv", "wmv", "mkv", "mpeg", "3gp", "ogv":
			typeFile = "video"
		default:
			log.Println("incorrect file's type")
			return "", errors.New("incorrect file type")
		}
	}
	return typeFile, nil
}

// getUserIdFromContext get user id from context
func getUserIdFromContext(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)

	if !ok {
		log.Println("no user id found")
		return 0, errors.New("user id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		log.Println("invalid user id type")
		return 0, errors.New("invalid user id type")
	}
	return idInt, nil
}

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}

/*
func hashPassword(pass string) string {
	hash := sha256.New()
	hash.Write([]byte(pass))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))

}


func (h *Handler) generateToken(name, pass string) (string, error) {
	user, err := h.repo.GetUser(name, hashPassword(pass))
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
*/
