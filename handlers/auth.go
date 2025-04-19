package handlers

import (
	"automation/db"
	"automation/models"
	"encoding/json"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if creds.Login == "" || creds.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	dbConn, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var user models.User
	err = dbConn.QueryRow(`
		SELECT workers.id, workers.login, workers.fio, positions.name as post, workers.pass 
		FROM workers 
		JOIN positions ON workers.post = positions.id 
		WHERE login = ?`,
		creds.Login).Scan(&user.Id, &user.Login, &user.Fio, &user.Post, &user.Pass)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !db.IsPasswordHashed(user.Pass) {
		hashedPassword, err := db.HashPassword(user.Pass)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		_, err = dbConn.Exec("UPDATE workers SET pass = ? WHERE login = ?", hashedPassword, user.Login)
		if err != nil {
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}
		user.Pass = hashedPassword
	}

	if !db.CheckPasswordHash(creds.Password, user.Pass) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	responseUser := struct {
		Id    int    `json:"id"`
		Login string `json:"login"`
		Fio   string `json:"fio"`
		Post  string `json:"post"`
	}{
		Id:    user.Id,
		Login: user.Login,
		Fio:   user.Fio,
		Post:  user.Post,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseUser)
}
