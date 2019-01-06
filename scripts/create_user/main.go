package main

import (
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/moonkeat/chainstack/services"
)

func main() {
	emailPtr := flag.String("email", "", "user email")
	passwordPtr := flag.String("password", "", "user password")
	isAdminPtr := flag.Bool("admin", false, "add admin user")

	flag.Parse()

	dbConnString := os.Getenv("DB_CONNSTRING")
	db, err := sqlx.Connect("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Failed to connect to postgres, connString: '%s'", dbConnString)
	}

	userService := services.NewUserService(db)
	err = userService.CreateUser(*emailPtr, *passwordPtr, *isAdminPtr)
	if err != nil {
		log.Fatal(err)
	}
}
