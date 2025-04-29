package handlers

import (
	"automation/db"     
	"automation/models" 
	"encoding/json"
	"log"
	"net/http"
)

func AddWorkerHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectDB() 
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer database.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var worker models.AddWorker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}
	hashedPassword, err := db.HashPassword(worker.Password)
	if err != nil{
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		log.Println("Password hash error:", err)
		return
	}
	query := `INSERT INTO workers (fio, post, login, pass)
VALUES (?, ?, ?, ?)
`
	_, err = database.Exec(query, worker.Name, worker.IdPosition, worker.Login, hashedPassword)
	if err != nil {
		http.Error(w, "Failed to add worker to database", http.StatusInternalServerError)
		log.Println("Error adding worker:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Worker added successfully"))
}
