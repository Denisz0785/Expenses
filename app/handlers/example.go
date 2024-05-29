package handlers

import "github.com/gin-gonic/gin"

type MySmallRepositoryForHandler interface {
	// здесь только те методы репозитория которые используются в ручке
}

type H func(c *gin.Context)

func NewMyHandler(repo MySmallRepositoryForHandler) H {
	return func(c *gin.Context) {

	}
}
