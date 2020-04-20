package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var (
	getProductsQuery   = "SELECT id, name, price FROM products LIMIT $1 OFFSET $2"
	getProductQuery    = "SELECT name, price FROM products WHERE id=$1"
	createProductQuery = "INSERT INTO products(name, price) VALUES($1, $2) RETURNING id"
	updateProductQuery = "UPDATE products SET name=$1, price=$2 WHERE id=$3"
	deleteProductQuery = "DELETE FROM products WHERE id=$1"
)

// product routes
func (a *App) productRoutes() {
	a.Router.HandleFunc("/products", a.getAllProducts).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product", a.createNewProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

func (a *App) getAllProducts(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getAllProductDetails(a.DB, getProductsQuery, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := fmt.Sprintf("Invalid product ID. Error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	p := product{ID: id}
	if err := p.getProductDetails(a.DB, getProductQuery); err != nil {
		switch err {
		case sql.ErrNoRows:
			msg := fmt.Sprintf("Product not found. Error: %s", err.Error())
			respondWithError(w, http.StatusNotFound, msg)
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) createNewProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p product
	if err := decoder.Decode(&p); err != nil {
		msg := fmt.Sprintf("Invalid request payload. Error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	defer r.Body.Close()

	if err := p.performCreateNewProduct(a.DB, createProductQuery); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := fmt.Sprintf("Invalid product ID. Error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	p := product{ID: id}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		msg := fmt.Sprintf("Invalid request payload. Error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	defer r.Body.Close()

	if err := p.performUpdateProduct(a.DB, updateProductQuery); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := fmt.Sprintf("Invalid product ID. Error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	p := product{ID: id}

	if err := p.performDeleteProduct(a.DB, deleteProductQuery); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
}
