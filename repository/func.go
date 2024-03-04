package repository

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// GetExpenseType gets one row of type of expenses from DB by name
func GetExpenseType(conn *pgx.Conn, name string) *string {

	var type_exp string

	err := conn.QueryRow(context.Background(),
		"SELECT type_expenses from expense_type, users where expense_type.users_id=users.id and users.name=$1", name).Scan(&type_exp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	return &type_exp
}

// GetManyRows gets all rows of type of expenses from DB by name
func GetManyRows(conn *pgx.Conn, name string) []string {
	rows, _ := conn.Query(context.Background(), "SELECT type_expenses from expense_type, users where expense_type.users_id=users.id and users.name=$1", name)
	numbers, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get rows: %v\n", err)
		os.Exit(1)
	}
	return numbers

}

// AddValuesDB insert row to the table user
func AddValuesDB(conn *pgx.Conn) error {

	commandTag, err := conn.Exec(context.Background(), "Insert into users(name,surname,login,pass,email) VALUES ('kolya','Bon','spiman','1243','er1@23.ru')")
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("add is not done")
	}
	fmt.Println("Adding row to the table users is finished")
	return err
}

// ConnDB connects to DB
func ConnDB(myurl string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv(myurl))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}
