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

	rows, err := dbConn.Query("SELECT w.id, w.login, w.fio, p.name AS post, w.pass FROM workers w LEFT JOIN positions p ON w.post = p.id WHERE w.fio LIKE ?",
	 "%"+name+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var workers []models.User
	for rows.Next() {
		var worker models.User
		if err := rows.Scan(&worker.Id, &worker.Login, &worker.Fio, &worker.Post, &worker.Pass); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		workers = append(workers, worker)
	}

	json.NewEncoder(w).Encode(workers)
}

//http://localhost:8080/get-worker-by-name/worker?name=Michael Johnson
