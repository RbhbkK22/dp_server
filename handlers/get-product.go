package handlers

import (
	"automation/db"
	"automation/models"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Error connecting to database:", err) 
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	rows, err := dbConn.Query("SELECT product.id, product.name, product.photo, product.discript as description, categories.name as category, brands.name as brand, product.quality as quantity, product.price FROM product JOIN categories on product.idCategories =categories.id JOIN brands on product.idBrands = brands.id")
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error retrieving products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Photo, &product.Description, &product.Category, &product.Brand, &product.Quantity, &product.Price); err != nil {
			log.Println("Error scanning product data:", err) 
			return
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error reading rows:", err) 
		http.Error(w, "Error reading products data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func GetProductImageHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Error connecting to database:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var photoPath string
	err = dbConn.QueryRow("SELECT photo FROM product WHERE id = ?", id).Scan(&photoPath)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
	
		return
	}
	baseDir := "E:\\back\\testv4hsserv\\automation\\uploads"

	imagePath := filepath.Join(baseDir, photoPath)
	log.Println("Trying to read image from path:", imagePath)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Println("Image file does not exist at path:", imagePath)
		http.Error(w, "Image not found", http.StatusNotFound)

		return
	}

	file, err := os.Open(imagePath)
	if err != nil {
		log.Println("Error opening image file:", err)
		http.Error(w, "Error opening image", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeFile(w, r, imagePath)
}
