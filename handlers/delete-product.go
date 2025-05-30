package handlers

import (
	"automation/db" 
	"encoding/json"
	"log"
	"net/http"
)

type DeleteProductRequest struct {
	ID int `json:"id"` 
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	database, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer database.Close()

	var req DeleteProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		log.Println("Error decoding request:", err)
		return
	}

	query := `DELETE FROM product WHERE id = ?`
	result, err := database.Exec(query, req.ID)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		log.Println("Error deleting product:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve delete result", http.StatusInternalServerError)
		log.Println("Error retrieving rows affected:", err)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No product found with the given criteria", http.StatusNotFound)
		log.Println("No matching product found")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product deleted successfully"))
}
