package handlers

import (
	"automation/db"
	"automation/models"
	"encoding/json"
	"log"
	"net/http"
)

func GetBrand(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Error connecting to database:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()
	rows, err := dbConn.Query("SELECT * FROM brands WHERE name LIKE ?", "%"+name+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var brands []models.Brand
	for rows.Next() {
		var brand models.Brand
		if err := rows.Scan(&brand.Id, &brand.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		brands = append(brands, brand)
	}
	json.NewEncoder(w).Encode(brands)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Error connecting to database:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()
	rows, err := dbConn.Query("SELECT * FROM categories WHERE name LIKE ?", "%"+name+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.Id, &category.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}
	json.NewEncoder(w).Encode(categories)
}

func GetPosition(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Error connecting to database:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()
	rows, err := dbConn.Query("SELECT * FROM positions WHERE name LIKE ?", "%"+name+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var positions []models.Position
	for rows.Next() {
		var position models.Position
		if err := rows.Scan(&position.Id, &position.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		positions = append(positions, position)
	}
	json.NewEncoder(w).Encode(positions)
}
