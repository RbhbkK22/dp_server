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

	rows, err := dbConn.Query("SELECT orders.id, clients.name AS client_name, orders.comment, orders.datatime, manager.fio AS manager_name, collector.fio AS collector_name, orders.status FROM orders JOIN clients ON orders.idClients = clients.id JOIN workers AS manager ON orders.idManager = manager.id LEFT JOIN workers AS collector ON orders.idCollector = collector.id;")
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
