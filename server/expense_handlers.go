package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type RepoExpense struct {
	Conn *pgx.Conn
}

type Expenses struct {
	Id    int64
	Title string
}

func (r *RepoExpense) GetTypesExpenseUser(ctx context.Context, id1 int) ([]Expenses, error) {
	rows, _ := r.Conn.Query(ctx, "SELECT id, title from expense_type where expense_type.users_id=$1", id1)
	expense, err := pgx.CollectRows(rows, pgx.RowToStructByName[Expenses])
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		return nil, err
	}
	return expense, nil
}

func (r *RepoExpense) GetExpenseHandler(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	if req.Method == http.MethodPost {
		id, err := strconv.Atoi(req.FormValue("id"))
		if err != nil {
			fmt.Fprintf(w, "Uncorrect value of id%v", err)
		}
		// []titleExpense keep title of expenses of  user
		titleExpense, err := r.GetTypesExpenseUser(ctx, id)
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
