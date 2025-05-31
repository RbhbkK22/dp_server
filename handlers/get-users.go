package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"automation/db"
	"automation/models"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Database connection error:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	rows, err := dbConn.Query(`
		SELECT w.id, w.login, w.fio, p.name AS post, w.pass
		FROM workers w
		LEFT JOIN positions p ON w.post = p.id
	`)

	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Login, &user.Fio, &user.Post, &user.Pass)
		if err != nil {
			log.Println("Error scanning user:", err)
			http.Error(w, "Error scanning user data", http.StatusInternalServerError)
			return
		}

		user.Pass = "" 
		users = append(users, user)
	}

	if len(users) == 0 {
		http.Error(w, "No users found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Println("Error encoding users:", err)
		http.Error(w, "Error encoding users data", http.StatusInternalServerError)
		return
	}
}
func GetClientsHandler(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	rows, err := dbConn.Query("SELECT id, name, contact FROM clients")
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var client models.Client
		err := rows.Scan(&client.ID, &client.Name, &client.Contact) 
		if err != nil {
			http.Error(w, "Error scanning client data", http.StatusInternalServerError)
			return
		}
		clients = append(clients, client)
	}

	if len(clients) == 0 {
		http.Error(w, "No clients found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(clients)
	if err != nil {
		http.Error(w, "Error encoding clients data", http.StatusInternalServerError)
		return
	}
}

func AddClient(w http.ResponseWriter, r *http.Request) {
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

	var client models.AddClientModel

	if err := json.NewDecoder(r.Body).Decode(&client); err !=nil{
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO clients (name, contact) VALUES (?, ?)`
	_, err = database.Exec(query, client.Name, client.Contact)
	if err != nil {
		http.Error(w, "Failed to add worker to database", http.StatusInternalServerError)
		log.Println("Error adding worker:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Worker added successfully"))
}
