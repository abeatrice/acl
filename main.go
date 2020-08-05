package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

type User struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
}

var DB *sql.DB
var err error

func index(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()

	selectStmt := `SELECT id, username, first_name, last_name, email FROM users WHERE 1=1 `

	if first_name, exists := params["first_name"]; exists {
		selectStmt += `AND UPPER(first_name) LIKE UPPER('%` + first_name[0] + `%')`
	}
	if last_name, exists := params["last_name"]; exists {
		selectStmt += `AND UPPER(last_name) LIKE UPPER('%` + last_name[0] + `%')`
	}
	if userName, exists := params["username"]; exists {
		selectStmt += `AND UPPER(username) LIKE UPPER('%` + userName[0] + `%')`
	}
	if email, exists := params["email"]; exists {
		selectStmt += `AND UPPER(email) LIKE UPPER('%` + email[0] + `%')`
	}

	stmt, err := DB.Prepare(selectStmt)
	check(err)
	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("users not found"))
		return
	}

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

func show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var user User

	row := DB.QueryRow("SELECT id, username, first_name, last_name, email FROM users where id = ?", vars["id"])
	err := row.Scan(&user.ID, &user.Username, &user.First_name, &user.Last_name, &user.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	json, err := json.Marshal(user)
	check(err)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(json))
}

func store(w http.ResponseWriter, r *http.Request) {

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == `` || user.First_name == `` || user.Last_name == `` || user.Email == `` {
		http.Error(w, `username, first_name, last_name, and email are required`, http.StatusBadRequest)
		return
	}

	stmt, err := DB.Prepare("INSERT INTO users (username, first_name, last_name, email) VALUES(?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.First_name, user.Last_name, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("user created"))
}

func update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("ok"))
}

func destroy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("ok"))
}

func main() {
	DB, err = sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DNS"))
	check(err)
	defer DB.Close()

	router := mux.NewRouter()

	router.HandleFunc("/users", index).Methods("GET")
	router.HandleFunc("/users/{id}", show).Methods("GET")
	router.HandleFunc("/users", store).Methods("POST")
	router.HandleFunc("/users/{id}", update).Methods("PUT", "PATCH")
	router.HandleFunc("/users/{id}", destroy).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
