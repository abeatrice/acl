package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type user struct {
	id         int
	username   string
	first_name string
	last_name  string
	email      string
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbDatabase := os.Getenv("DB_DATABASE")

	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbDatabase) //user:pass@tcp(host)/database

	db, err := sql.Open(os.Getenv("DB_DRIVER"), dns)
	check(err)
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, username, first_name, last_name, email FROM users")
	check(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	check(err)

	var (
		id         int
		username   string
		first_name string
		last_name  string
		email      string
	)

	for rows.Next() {
		err := rows.Scan(&id, &username, &first_name, &last_name, &email)
		check(err)
		fmt.Printf("%T: %v\n", username, username)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Ok"))
}

func main() {
	http.HandleFunc("/users", userHandler)
	http.ListenAndServe(":8000", nil)
}

func check(err error) {
	if err != nil {
		// log.Fatal(err)
		panic(err)
	}
}
