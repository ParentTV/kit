package infrastructure

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"kit/auth"
	"kit/event_bus"
	"log"
	"os"
)

func Setup() (*sql.DB, *event_bus.EventBus, *auth.Auth) {
	var (
		dbUser      = os.Getenv("DB_USER")
		dbPass      = os.Getenv("DB_PASSWORD")
		dbHost      = os.Getenv("DB_HOST")
		dbPort      = os.Getenv("DB_PORT")
		dbName      = os.Getenv("DB_NAME")
		serviceName = os.Getenv("SERVICE_NAME")
	)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)
	fmt.Println(dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	eb := event_bus.NewEventBus(serviceName)
	return db, eb, auth.NewAuth(eb)
}
