package main

import (
	"flag"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/moonkeat/chainstack/services"
)

func main() {
	emailPtr := flag.String("email", "", "user email")
	passwordPtr := flag.String("password", "", "user password")
	isAdminPtr := flag.Bool("admin", false, "add admin user")
	quotaPtr := flag.Int("quota", -1, "user quota")

	flag.Parse()

	dbConnString := os.Getenv("DB_CONNSTRING")
	db, err := sqlx.Connect("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Failed to connect to postgres, connString: '%s'", dbConnString)
	}

	var quota *int
	if *quotaPtr != -1 {
		quota = quotaPtr
	}

	userService := services.NewUserService(db)
	_, err = userService.CreateUser(*emailPtr, *passwordPtr, *isAdminPtr, quota)
	if err != nil {
		log.Fatal(err)
	}
}
