package main

import (
	"berniestack/internal/features/subject"
	"database/sql"
	"fmt"
	"log"
	"os"

	data "berniestack/internal/db"

	"github.com/go-fuego/fuego"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	s := fuego.NewServer()
	port := os.Getenv("PORT")
	connStr := os.Getenv("DATABASE_URL")
	fuego.WithAddr("localhost:" + port)(s)

	conn, err := sql.Open("sqlite3", connStr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	db := data.New(conn)
	rs := subject.NewSubjectResources(db)
	rs.MountRoutes(s)

	err = s.Run()
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
	log.Println("Shutdown...")
}
