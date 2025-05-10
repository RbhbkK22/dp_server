package handlers

import (
	"automation/db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, `{"error":"Only POST method allowed"}`, http.StatusMethodNotAllowed)
        return
    }

    db, err := db.ConnectDB()
    if err != nil {
        http.Error(w, `{"error":"Database connection failed"}`, http.StatusInternalServerError)
        log.Printf("DB connection error: %v", err)
        return
    }
    defer db.Close()

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, `{"error":"Failed to read request body"}`, http.StatusBadRequest)
        return
    }

    var request struct {
        ClientID int            `json:"clientId"`
        Comment  string         `json:"comment"`
        Items    map[string]int `json:"items"` 
    }

    if err := json.Unmarshal(body, &request); err != nil {
        http.Error(w, `{"error":"Invalid JSON format"}`, http.StatusBadRequest)
        return
    }

    var clientExists bool
    err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM clients WHERE id = ?)", request.ClientID).Scan(&clientExists)
    if err != nil || !clientExists {
        http.Error(w, `{"error":"Client not found"}`, http.StatusBadRequest)
        return
    }

    if len(request.Items) == 0 {
        http.Error(w, `{"error":"No items in order"}`, http.StatusBadRequest)
        return
    }

    tx, err := db.Begin()
    if err != nil {
        http.Error(w, `{"error":"Transaction start failed"}`, http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    res, err := tx.Exec(
        "INSERT INTO orders (idClients, comment, datatime, idManager, status) VALUES (?, ?, ?, 1, 'new')",
        request.ClientID, request.Comment, time.Now())
    if err != nil {
        http.Error(w, `{"error":"Failed to create order"}`, http.StatusInternalServerError)
        return
    }

    orderID, err := res.LastInsertId()
    if err != nil {
        http.Error(w, `{"error":"Failed to get order ID"}`, http.StatusInternalServerError)
        return
    }

    for strID, quantity := range request.Items {
        itemID, err := strconv.Atoi(strID)
        if err != nil {
            http.Error(w, `{"error":"Invalid item ID format"}`, http.StatusBadRequest)
            return
        }

        var stock int
        err = tx.QueryRow("SELECT quality FROM product WHERE id = ?", itemID).Scan(&stock)
        if err != nil {
            http.Error(w, fmt.Sprintf(`{"error":"Product %d not found"}`, itemID), http.StatusBadRequest)
            return
        }

        if quantity <= 0 || quantity > stock {
            http.Error(w, fmt.Sprintf(`{"error":"Invalid quantity for product %d"}`, itemID), http.StatusBadRequest)
            return
        }

        _, err = tx.Exec(
            "INSERT INTO items (idItems, quality, idOrders) VALUES (?, ?, ?)",
            itemID, quantity, orderID)
        if err != nil {
            http.Error(w, `{"error":"Failed to add items"}`, http.StatusInternalServerError)
            return
        }

        _, err = tx.Exec(
            "UPDATE product SET quality = quality - ? WHERE id = ?",
            quantity, itemID)
        if err != nil {
            http.Error(w, `{"error":"Failed to update stock"}`, http.StatusInternalServerError)
            return
        }
    }

    if err := tx.Commit(); err != nil {
        http.Error(w, `{"error":"Failed to complete order"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "orderId": orderID,
    })
}