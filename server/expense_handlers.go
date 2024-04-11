package server

import (
	"context"
	"encoding/json"
	dto "expenses/dto_expenses"
	"expenses/repository"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Server struct {
	repo *repository.ExpenseRepo
}

func NewServer(c *repository.ExpenseRepo) *Server {
	return &Server{repo: c}
}

func (r *Server) GetExpenseHandler(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	user := &dto.User{}
	if req.Method == http.MethodPost {
		if checkId := req.FormValue("id"); checkId != "" {
			id, err := strconv.Atoi(checkId)
			if err != nil {
				fmt.Fprintf(w, "Uncorrect value of id%v", err)
			}
			user.Id = id
		}
		login := req.FormValue("login")
		user.Login = login
		name := req.FormValue("name")
		user.Name = name
		fmt.Println(user)
		// []titleExpense keep title of expenses of  user
		titleExpense, err := r.repo.GetTypesExpenseUser(ctx, user)
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
