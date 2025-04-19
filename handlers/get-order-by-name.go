package handlers

import (
	"automation/db"
	"automation/models"
	"encoding/json"
	"net/http"
)

func GetOrdersByClientName(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.ConnectDB()
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	rows, err := dbConn.Query(`
SELECT orders.id, clients.name, orders.comment, orders.datatime,
       manager.fio AS manager_name,
       collector.fio AS collector_name,
       orders.status
FROM orders
INNER JOIN clients ON orders.idClients = clients.id
INNER JOIN workers AS manager ON orders.idManager = manager.id
LEFT JOIN workers AS collector ON orders.idCollector = collector.id
WHERE clients.name LIKE ?
`, "%"+name+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.Id, &order.ClientName, &order.Comment, &order.Date,
			&order.ManagerName, &order.CollectorName, &order.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	json.NewEncoder(w).Encode(orders)
}

//http://localhost:8080/get-order-by-name/order?name=John Doe
