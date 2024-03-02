package main

import (
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
	var s Conn_DB
	err := envconfig.Process("PGEXP", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	format := "host: %s\nuser: %s\npassword: %s\ndatabase: %s\n"
	_, err = fmt.Printf(format, s.Host, s.User, s.Password, s.DbName)
	if err != nil {
		log.Fatal(err.Error())
	}

	type_expenses := repository.Get_expense_types("Ivan")
	fmt.Println(*type_expenses)
	fmt.Println()
	// ./expenses cmd=get_expense_types user=vasya
	funcPtr := flag.String("cmd", "none", "function")
	userPtr := flag.String("name", "none", "user's name")
	flag.Parse()
	fmt.Println("function", *funcPtr)
	fmt.Println("user's name", *userPtr)
	fmt.Println("tail", flag.Args())

	var resultExpenses *string
	if strings.EqualFold(*funcPtr, "Get_expense_types") {
		resultExpenses = repository.Get_expense_types(*userPtr)
	}
	fmt.Printf("Expenses_type of %s = %s\n", *userPtr, *resultExpenses)

}
