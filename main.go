package main

import (
	"context"
	"expenses/getenv"
	"expenses/repository"
	"flag"
	"fmt"
	"strings"
)

func main() {
	// GetEnvVarForDB() gets environment variables by prefix, adds to struct and prints it
	getenv.GetEnvVarForDB()

	// ConnectToDB connects to DB
	myUrl := "MYURL"
	conn, err := repository.ConnectToDB(myUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close(context.Background())

	// Create new structure wich consist data about connection with database
	ConnExpRepo := repository.NewExpenseRepo(conn)

	//GetExpenseType gets one row of type of expenses from DB by name
	name := "Igor"
	typeExpenses, err1 := ConnExpRepo.GetExpenseType(name)
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	fmt.Printf("one type of expenses %s: %s\n", name, *typeExpenses)
	fmt.Println()

	//GetManyRows gets all rows of type of expenses from DB by name
	rows, err2 := ConnExpRepo.GetManyRowsByName(name)
	if err2 != nil {
		fmt.Println(err2.Error())
	}
	fmt.Printf("all types of expenses %s: %v\n", name, rows)
	fmt.Println()

	// AddValuesDB inserts a row to the table users
	err2 = ConnExpRepo.AddValuesDB()
	if err2 != nil {
		fmt.Println(err2.Error())
	}
	fmt.Println()

	// define flags for getting values of flags command ./expenses cmd=get_expense_types user=Ivan and
	// ./expenses -cmd=Add -login=igor23 -exp_type=swimming -time=2024-02-25-17:26 -spent=500
	funcPtr := flag.String("cmd", "none", "function")
	userPtr := flag.String("name", "none", "user's name")
	loginPtr := flag.String("login", "none", "user's login")
	expTypePtr := flag.String("exp_type", "none", "type of expenses")
	timePtr := flag.String("time", "none", "time of expenses")
	spentPtr := flag.Float64("spent", 0.0, "amount of expenses")
	//Parse() parses the command line into the defined flags
	flag.Parse()

	// define which command was input
	switch {
	case strings.EqualFold(*funcPtr, "Get_ManyRows"):
		var resultExpenses []string
		resultExpenses, err = ConnExpRepo.GetManyRowsByName(*userPtr)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println()
		fmt.Printf("Expenses_type of %s = %s\n", *userPtr, resultExpenses)

	case strings.EqualFold(*funcPtr, "add"):
		//
		err := ConnExpRepo.AddExpTransaction(loginPtr, expTypePtr, timePtr, spentPtr)
		if err != nil {
			fmt.Println(err.Error())
		}
	default:
		fmt.Println("check your input data in a command-line")
	}

}
