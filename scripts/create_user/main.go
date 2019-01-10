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
	quotaPtr := flag.Int("quota", services.UserQuotaUndefined, "user quota")

	flag.Parse()

	dbConnString := os.Getenv("DB_CONNSTRING")
	db, err := sqlx.Connect("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Failed to connect to postgres, connString: '%s'", dbConnString)
	}

	var quota *int
	if *quotaPtr != services.UserQuotaUndefined {
		quota = quotaPtr
	}
	if quota != nil && *quota < 0 {
		log.Fatalf("User quota should be -1 (unlimited quota) or at least 0")
	}

	userService := services.NewUserService(db)
	_, err = userService.CreateUser(*emailPtr, *passwordPtr, *isAdminPtr, quota)
	if err != nil {
		log.Fatal(err)
	}
}
