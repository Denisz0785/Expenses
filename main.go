// main package of project
package main

import (
	"context"
	"expenses/repository"
	"expenses/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Println("Не удалось открыть файл логов, используется стандартный stderr")
	}
	defer file.Close()

	// ConnectToDB connects to DB
	myURL := "MYURL"
	ctx := context.Background()
	conn, err := repository.ConnectToDB(ctx, myURL)
	if err != nil {
		log.Fatalf("error connect DB:%s", err.Error())

	}
	defer conn.Close(ctx)
	// initializing config from file
	if err := initConfig(); err != nil {
		log.Fatalf("error of initializing configs:%s", err.Error())
	}

	/*
		// define flags for getting values of flags command ./expenses cmd=get_expense_types user=Ivan and
		// ./expenses -cmd=Add -login=igor23 -exp_type=swimming -time=2024-02-25-17:26 -spent=500
		funcPtr := flag.String("cmd", "none", "function")
		userPtr := flag.String("name", "none", "user's name")
		loginPtr := flag.String("login", "none", "user's login")
		expTypePtr := flag.String("exp_type", "none", "type of expenses")
		timePtr := flag.String("time", "none", "time of expenses")
		spentPtr := flag.Float64("spent", 0.0, "amount of expenses")

		//Parse parses the command line into the defined flags
		flag.Parse()

	*/

	// define which command was input
	/*
		switch {

		case strings.EqualFold(*funcPtr, "Get_ManyRows"):
			user := &dto.TypesExpenseUserParams{Name: *userPtr}
			resultExpenses, err := repos.GetTypesExpenseUser(ctx, user)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println()
			fmt.Printf("Expenses_type of %v = %v\n", user.Name, resultExpenses)

		case strings.EqualFold(*funcPtr, "add"):

			err := repos.CreateUserExpense(ctx, loginPtr, expTypePtr, timePtr, spentPtr)
			if err != nil {
				fmt.Println(err.Error())
			}
		case strings.EqualFold(*funcPtr, "run_server"):

	*/
	// Create new structure wich consist data about connection with database
	repos := repository.NewExpenseRepo(conn)
	handlers := server.NewHandler(repos)

	srv := new(server.Server)

	// create signal channel
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	// run server in a goroutine
	go func() {
		if err := srv.Run(handlers.InitRoutes(), viper.GetString("port")); err != nil {
			log.Fatalf("error run server: %s", err.Error())
		}
	}()
	<-ch
	log.Println("server shutting down")
	// Shutdown gracefully shuts down the server without interrupting any active connections
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error of shutdown server:%s", err.Error())

	}
}

// initConfig initializes configs from file
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()

}
