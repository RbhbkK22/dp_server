package handlers

import (
	"automation/db"
	"automation/models"
	"encoding/json"
	"log"
	"net/http"
)

func GetOrders(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Database connection error:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	rows, err := dbConn.Query("SELECT orders.id, clients.name AS client_name, orders.comment, orders.datatime, manager.fio AS manager_name, collector.fio AS collector_name, orders.status FROM orders JOIN clients ON orders.idClients = clients.id JOIN workers AS manager ON orders.idManager = manager.id LEFT JOIN workers AS collector ON orders.idCollector = collector.id ORDER BY orders.datatime DESC;")
	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var orders []models.Order

	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.Id, &order.ClientName, &order.Comment, &order.Date, &order.ManagerName, &order.CollectorName, &order.Status)
		if err != nil {
			log.Println("Database scan error:", err)
			http.Error(w, "Database scan error", http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		log.Println("Encode error:", err)
		http.Error(w, "Encode error", http.StatusInternalServerError)
		return
	}

}

func GetItemsInOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id parameter dont faund", http.StatusBadRequest)
	}
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Database connection error:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	rows, err := dbConn.Query(`SELECT items.idItems, product.name, product.photo, items.quality
		FROM items
		LEFT JOIN product ON product.id = items.idItems
		WHERE items.idOrders = ?;`, id)
	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var result []models.ItemInfo

	for rows.Next() {
		var item models.ItemInfo
		err := rows.Scan(&item.ItemId, &item.Name, &item.Photo, &item.Quality)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, item)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Encode error:", err)
		http.Error(w, "Encode error", http.StatusInternalServerError)
	}
}

// http://localhost:8080/get-items-in-order?id=1
