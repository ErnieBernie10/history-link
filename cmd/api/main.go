package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	data "historylink/internal/db"
	"historylink/internal/features/record"

	"github.com/go-fuego/fuego"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {

	port := os.Getenv("PORT")
	connStr := os.Getenv("DATABASE_URL")
	s := fuego.NewServer(
		fuego.WithAddr("localhost:" + port),
	)

	conn, err := sql.Open("postgres", connStr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	db := data.New(conn)
	rs := record.NewRecordResources(db)
	rs.MountRoutes(s)

	err = s.Run()
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
	log.Println("Shutdown...")
}
