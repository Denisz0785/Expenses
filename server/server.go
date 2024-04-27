package server

import (
	"context"
	"net/http"
	"time"
)

type Connect struct {
	httpServer *http.Server
}

// Run create router and run a server
func (c *Connect) Run(r *Server, port string) {
	mux := http.NewServeMux()
	c.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	mux.HandleFunc("/expense/list/", r.GetExpenseHandler)
	mux.HandleFunc("/expense/upload/", r.UploadFile)
	mux.HandleFunc("/expense/delete/", r.DeleteFile)
	c.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any active connections
func (c *Connect) Shutdown(ctx context.Context) error {
	return c.httpServer.Shutdown(ctx)
}
