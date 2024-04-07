package server

import (
	"net/http"
)

func Run(r *RepoExpense) {
	mux := http.NewServeMux()
	mux.HandleFunc("/expense/list/", r.GetExpenseHandler)
	http.ListenAndServe(":8080", mux)
}
