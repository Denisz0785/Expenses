package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func Get_expense_types(name string) *string {

	conn, err := pgx.Connect(context.Background(), os.Getenv("MYURL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())
	var type_exp string

	err = conn.QueryRow(context.Background(),
		"SELECT type_expenses from expense_type, users where expense_type.users_id=users.id and users.name=$1", name).Scan(&type_exp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	return &type_exp
}
