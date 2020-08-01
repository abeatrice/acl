package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
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
	var users []User

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.First_name, &user.Last_name, &user.Email)
		check(err)

		users = append(users, user)
	}

	json, err := json.Marshal(users)
	check(err)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(json))
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbDatabase := os.Getenv("DB_DATABASE")
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbDatabase) //user:pass@tcp(host)/database

	db, err := sql.Open(dbDriver, dns)
	check(err)
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	var user User

	row := db.QueryRow("SELECT id, username, first_name, last_name, email FROM users where id = ?", id)
	err = row.Scan(&user.ID, &user.Username, &user.First_name, &user.Last_name, &user.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	json, err := json.Marshal(user)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(json))
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	router.HandleFunc("/users", usersHandler)
	router.HandleFunc("/user/{id}", userHandler)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
