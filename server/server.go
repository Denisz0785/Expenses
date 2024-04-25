package server

import (
	"net/http"
)

func Run(r *Server) {
	mux := http.NewServeMux()
	mux.HandleFunc("/expense/list/", r.GetExpenseHandler)
	mux.HandleFunc("/expense/upload/", r.UploadFile)
	http.ListenAndServe(":8080", mux)
}
