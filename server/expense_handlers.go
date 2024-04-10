package server

import (
	"context"
	"encoding/json"
	"expenses/repository"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Server struct {
	conn *repository.ExpenseRepo
}

func NewServer(c *repository.ExpenseRepo) *Server {
	return &Server{conn: c}
}

func (r *Server) GetExpenseHandler(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	if req.Method == http.MethodPost {
		id, err := strconv.Atoi(req.FormValue("id"))
		if err != nil {
			fmt.Fprintf(w, "Uncorrect value of id%v", err)
		}
		// []titleExpense keep title of expenses of  user
		titleExpense, err := r.conn.GetTypesExpenseUser(ctx, id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// write to w json data of expenseList
		err = json.NewEncoder(w).Encode(titleExpense)
		if err != nil {
			fmt.Println("Error of marshalig to json")
		}
		w.Header().Set("Content-Type", "application/json")
		return
	} else {
		io.WriteString(w, "enter a post method")
		return
	}
}
