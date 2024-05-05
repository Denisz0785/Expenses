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
	myUrl := "MYURL"
	ctx := context.Background()
	conn, err := repository.ConnectToDB(ctx, myUrl)
	if err != nil {
		log.Fatalf("error connect DB:%s", err.Error())

	}
	defer conn.Close(ctx)
	// initializing config from file
	if err := initConfig(); err != nil {
		log.Fatalf("error of initializing configs:%s", err.Error())
	}
	// Create new structure wich consist data about connection with database
	ConnExpRepo := repository.NewExpenseRepo(conn)
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
			resultExpenses, err := ConnExpRepo.GetTypesExpenseUser(ctx, user)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println()
			fmt.Printf("Expenses_type of %v = %v\n", user.Name, resultExpenses)

		case strings.EqualFold(*funcPtr, "add"):

			err := ConnExpRepo.CreateUserExpense(ctx, loginPtr, expTypePtr, timePtr, spentPtr)
			if err != nil {
				fmt.Println(err.Error())
			}
		case strings.EqualFold(*funcPtr, "run_server"):

	*/
	str := server.NewServer(ConnExpRepo)
	srv := new(server.Connect)
	// create signal channel
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	// run server in a goroutine
	go func() {
		if err := srv.Run(str.InitRoutes(), viper.GetString("port")); err != nil {
			log.Fatalf("error run server: %s", err.Error())
		}
	}()
	<-ch
	log.Println("server shutting down")
	// Shutdown gracefully shuts down the server without interrupting any active connections
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error of shutdown server:%s", err.Error())
		/*	}

			default:
				fmt.Println("check your input data in a command-line")
			}

		*/

	}
}

// initConfig initializes configs from file
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()

}
