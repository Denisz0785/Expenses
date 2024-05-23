package server

import (
	dto "expenses/dto_expenses"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	userCtx = "userId"
)

// signUp registers new user
func (h *Handler) signUp(c *gin.Context) {
	var user dto.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("incorrect user data:%s", err.Error())
		newErrorResponse(c, http.StatusBadRequest, "incorrect user data")
		return
	}
	user.Pass = hashPassword(user.Pass)
	id, err := h.repo.CreateUser(c, &user)
	if err != nil {
		log.Printf("error create user:%s", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, "error of create user")
		return
	}
	c.JSON(http.StatusOK, map[string]int{
		"id": id,
	})
}

// input save user data
type input struct {
	Name string `json:"name" binding:"required"`
	Pass string `json:"pass" binding:"required"`
}

// signIn authenticate user
func (h *Handler) signIn(c *gin.Context) {
	var user input
	err := c.BindJSON(&user)
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.generateToken(user.Name, user.Pass)
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

// userIdentity define id user by token and save id user to context
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		log.Println("empty authorization header")
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		log.Println("incorrect authorization header")
		newErrorResponse(c, http.StatusUnauthorized, "incorrect auth header")
		return
	}
	if len(headerParts[1]) == 0 {
		log.Println("missing token")
		newErrorResponse(c, http.StatusUnauthorized, "missed token")
		return
	}
	userId, err := parseToken(headerParts[1])
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
	c.Set(userCtx, userId)
}
