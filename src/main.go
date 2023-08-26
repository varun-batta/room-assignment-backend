package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type user struct {
	Name     string
	Username string
	Password string
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signUpRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

var db *sql.DB

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var body loginRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Query database for user
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE username='%s' AND password='%s'", body.Username, body.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users := []user{}
	for rows.Next() {
		var u user
		err = rows.Scan(&u.Name, &u.Username, &u.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	if len(users) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Handle the response
	responseData := authResponse{
		Token: "test1234",
		Name:  users[0].Name,
	}
	response, err := json.Marshal(responseData)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var body signUpRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add user to database and sign them in
	sqlStatement := fmt.Sprintf("INSERT INTO users VALUES ('%s', '%s', '%s')", body.Name, body.Username, body.Password)
	_, err = db.Exec(sqlStatement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle the response
	responseData := authResponse{
		Token: "test1234",
		Name:  body.Name,
	}
	response, err := json.Marshal(responseData)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func main() {
	connStr := "user=admin dbname=roomassignment sslmode=disable"
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db = database

	// Backend Code
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/signup", signup)
	log.Fatal(http.ListenAndServe(":7060", nil))
}
