package handlers

import (
	"automation/db"
	"automation/models"
	"encoding/json"
	"net/http"
	"strings"
)

func GetProductByName(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing product name", http.StatusBadRequest)
		return
	}

	name = strings.Trim(name, "'")

	query := `
		SELECT 
			p.id,
			p.name,
			p.price,
			p.photo,
			c.name AS category,
			b.name AS brand,
			p.quality AS quantity,
			p.discript AS description
		FROM product p
		JOIN categories c ON p.idCategories = c.id
		JOIN brands b ON p.idBrands = b.id
		WHERE p.name LIKE ?;
	`

	rows, err := dbConn.Query(query, "%"+name+"%")
	if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Photo, &p.Category, &p.Brand, &p.Quantity, &p.Description)
		if err != nil {
			http.Error(w, "Failed to scan product: "+err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

//http://localhost:8080/get-product-by-name/product?name=название

type ProductName struct {
	Name string `json:"name"`
}

func GetProductNameById(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	productId := r.URL.Query().Get("id")
	if productId == "" {
		http.Error(w, "Missing product id", http.StatusBadRequest)
		return
	}

	productId = strings.Trim(productId, "'") // это можно убрать, если параметр приходит корректно

	query := `SELECT name FROM product WHERE id = ?`

	var product ProductName
	err = dbConn.QueryRow(query, productId).Scan(&product.Name)
	if err != nil {
		http.Error(w, "Query error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

//http://localhost:8080/get-product-name-by-id/product?id=id
