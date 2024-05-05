package server

import (
	dto "expenses/dto_expenses"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// signUp registers new user
func (s *Server) signUp(c *gin.Context) {
	var user dto.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("incorrect user data:%s", err.Error())
		c.JSON(http.StatusBadRequest, "incorrect user data")
		return
	}
	id, err := s.repo.CreateUser(c, &user)
	if err != nil {
		log.Printf("error create user:%s", err.Error())
		c.JSON(http.StatusInternalServerError, "error of create user")
		return
	}
	c.JSON(http.StatusOK, map[string]int{
		"id": id,
	})

}
