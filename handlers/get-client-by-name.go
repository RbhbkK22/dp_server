package handlers

import (
	"automation/db"
	"automation/models"
	"net/http"
	"encoding/json"
)

func GetClientsByName(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.ConnectDB()
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	rows, err := dbConn.Query("SELECT id, name, contact FROM clients WHERE name LIKE ?", "%"+name+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var client models.Client
		if err := rows.Scan(&client.ID, &client.Name, &client.Contact); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clients = append(clients, client)
	}

	json.NewEncoder(w).Encode(clients)
}
//GET http://localhost:8080/get-client-by-name/client?name=John Doe