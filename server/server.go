package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Connect struct {
	httpServer *http.Server
}

func (s *Server) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", s.signUp)
		//		auth.POST("/sign-in", s.signIn)
	}
	api := router.Group("/api") //s.userIdentity
	{
		expenses := api.Group("/expenses")
		{
			expenses.GET("/", s.GetExpenseHandler)
		}
		files := api.Group("/files")
		{
			files.POST("/:id", s.UploadFile)
			files.DELETE("/:id", s.DeleteFile)
		}
	}
	return router
}

// Run create router and run a server
func (c *Connect) Run(handler http.Handler, port string) error {
	//mux := http.NewServeMux()
	c.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	//mux.HandleFunc("/expense/list/", r.GetExpenseHandler)
	//mux.HandleFunc("/expense/upload/", r.UploadFile)
	//mux.HandleFunc("/expense/delete/", r.DeleteFile)
	return c.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any active connections
func (c *Connect) Shutdown(ctx context.Context) error {
	return c.httpServer.Shutdown(ctx)
}
