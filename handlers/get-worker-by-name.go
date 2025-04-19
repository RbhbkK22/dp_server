package handlers

import (
	"automation/db"
	"automation/models"
	"encoding/json"
	"net/http"
)

func GetWorkersByName(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.ConnectDB()
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	rows, err := dbConn.Query("SELECT id, fio, login, pass FROM workers WHERE fio LIKE ?", "%"+name+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var workers []models.Worker
	for rows.Next() {
		var worker models.Worker
		if err := rows.Scan(&worker.Name, &worker.Position, &worker.Login, &worker.Password); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		workers = append(workers, worker)
	}

	json.NewEncoder(w).Encode(workers)
}

//http://localhost:8080/get-worker-by-name/worker?name=Michael Johnson
