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

type RepoExpense repository.ExpenseRepo

type Expenses struct {
	Id    int64
	Title string
}

func (r *RepoExpense) GetExpenseHandler(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	if req.Method == http.MethodPost {
		id, err := strconv.Atoi(req.FormValue("id"))
		if err != nil {
			fmt.Fprintf(w, "Uncorrect value of id%v", err)
		}
		// convert r to ExpenseRepo from repository to call methods from repository
		expenseRepo := repository.ExpenseRepo(*r)
		// []titleExpense keep title of expenses of  user
		titleExpense, err := expenseRepo.GetTypesExpenseUser(ctx, id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// []idExpenseString keep Id of expenses of user
		idExpenseString, err := expenseRepo.GetIdExpenseUser(ctx, id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		idExpenseInt := make([]int, len(idExpenseString))
		idExpenseInt64 := make([]int64, len(idExpenseInt))
		// convert string type of ID to Int
		for i, v := range idExpenseString {
			idExpenseInt[i], err = strconv.Atoi(v)
			if err != nil {
				fmt.Fprintf(w, "Uncorrect type of id%v", err)
			}
		}
		// convert int type of ID to int64
		for i, v := range idExpenseInt {
			idExpenseInt64[i] = int64(v)
		}
		// []expenseList keep list of id and title of user's type of expenses
		expenseList := make([]Expenses, 0, len(idExpenseInt64))
		for i := 0; i < len(titleExpense); i++ {
			expenseList = append(expenseList, Expenses{Id: idExpenseInt64[i], Title: titleExpense[i]})
		}
		fmt.Println(expenseList)
		// write to w json data of expenseList
		err = json.NewEncoder(w).Encode(expenseList)
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
