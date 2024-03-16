package getenv

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type EnvVarForDB struct {
	Host     string //`envconfig:"PGEXP_HOST"`
	User     string //`envconfig:"PGEXP_USER"`
	Password string //`envconfig:"PGEXP_PASSWORD"`
	DbName   string //`envconfig:"PGEXP_DBNAME"`
}

// GetEnvVarForDB() gets environment variables by prefix, adds to struct and prints it
func GetEnvVarForDB() error {
	var strEnvVar EnvVarForDB
	err := envconfig.Process("PGEXP", &strEnvVar)
	if err != nil {
		log.Fatal(err.Error())
	}

	format := "Date of environment variables: host: %s\nuser: %s\npassword: %s\ndatabase: %s\n"
	_, err = fmt.Printf(format, strEnvVar.Host, strEnvVar.User, strEnvVar.Password, strEnvVar.DbName)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println()
	return err
}
