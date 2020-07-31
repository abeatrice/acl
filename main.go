package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type User struct {
	ID         int
	Username   string
	First_name string
	Last_name  string
	Email      string
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbDatabase := os.Getenv("DB_DATABASE")
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbDatabase) //user:pass@tcp(host)/database

	db, err := sql.Open(dbDriver, dns)
	check(err)
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, username, first_name, last_name, email FROM users")
	check(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	check(err)

	var user User

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.First_name, &user.Last_name, &user.Email)
		check(err)

		// js, err := json.Marshal(user)
		// check(err)

		fmt.Printf("%T: %v\n", user, user)
		// fmt.Printf("%T: %v\n", js, js)
		// w.Header().Set("Content-Type", "application/json")
		// w.Write(js)
	}
	fmt.Println("here2")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Ok"))
}

func main() {
	http.HandleFunc("/users", usersHandler)
	http.ListenAndServe(":8000", nil)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
