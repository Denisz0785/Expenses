package main

import (
	"context"
	"expenses/repository"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Conn_DB struct {
	Host     string //`envconfig:"PGEXP_HOST"`
	User     string //`envconfig:"PGEXP_USER"`
	Password string //`envconfig:"PGEXP_PASSWORD"`
	DbName   string //`envconfig:"PGEXP_DBNAME"`
}

func main() {
	// add values to struct from environment variables by prefix
	var s Conn_DB
	err := envconfig.Process("PGEXP", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	format := "Date of environment variables: host: %s\nuser: %s\npassword: %s\ndatabase: %s\n"
	_, err = fmt.Printf(format, s.Host, s.User, s.Password, s.DbName)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println()

	// ConnDB connects to DB
	myUrl := "MYURL"
	conn := repository.ConnDB(myUrl)
	defer conn.Close(context.Background())

	//GetExpenseType gets one row of type of expenses from DB by name
	name := "Ivan"
	type_expenses := repository.GetExpenseType(conn, name)
	fmt.Printf("type of expenses %s: %s\n", name, *type_expenses)
	fmt.Println()

	//GetManyRows gets all rows of type of expenses from DB by name
	rows := repository.GetManyRows(conn, name)
	fmt.Printf("all types of expenses %s: %v\n", name, rows)
	fmt.Println()

	// AddValuesDB insert row to the table user
	repository.AddValuesDB(conn)
	fmt.Println()

	// create command ./expenses cmd=get_expense_types user=Ivan
	funcPtr := flag.String("cmd", "none", "function")
	userPtr := flag.String("name", "none", "user's name")
	flag.Parse()
	fmt.Println("Values of flags are:")
	fmt.Println("function:", *funcPtr)
	fmt.Println("user's name:", *userPtr)
	fmt.Println("tail:", flag.Args())

	var resultExpenses []string
	if strings.EqualFold(*funcPtr, "GetManyRows") {
		resultExpenses = repository.GetManyRows(conn, *userPtr)
	}
	fmt.Println()
	fmt.Printf("Expenses_type of %s = %s\n", *userPtr, resultExpenses)

}
