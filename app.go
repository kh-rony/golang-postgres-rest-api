package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

//struct to hold refs of router and database
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// create database connection and set up routing
func (a *App) Initialize(user, password, dbname, host, port, sslmode string) {

	dataSourceName := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", user, password, dbname, host, port, sslmode)

	var err error
	a.DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	err = a.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// run application
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))

	defer a.DB.Close()
}

// initialize routes into router that call methods on requests
func (a *App) initializeRoutes() {
	a.productRoutes()
}

// send a payload of JSON content
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// send a JSON error message
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
