package main

import (
	"database/sql"
)

// the product type is the fields that cover each item in our database
// struct fields include encoded JSON key names
type product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func getAllProductDetails(db *sql.DB, getProductsQuery string, start, count int) ([]product, error) {
	rows, err := db.Query(getProductsQuery, count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (p *product) getProductDetails(db *sql.DB, getProductQuery string) error {
	return db.QueryRow(getProductQuery, p.ID).Scan(&p.Name, &p.Price)
}

func (p *product) performCreateNewProduct(db *sql.DB, createProductQuery string) error {
	err := db.QueryRow(createProductQuery, p.Name, p.Price).Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *product) performUpdateProduct(db *sql.DB, updateProductQuery string) error {
	_, err := db.Exec(updateProductQuery, p.Name, p.Price, p.ID)
	return err
}

func (p *product) performDeleteProduct(db *sql.DB, deleteProductQuery string) error {
	_, err := db.Exec(deleteProductQuery, p.ID)
	return err
}
