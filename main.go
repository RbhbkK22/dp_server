package main

import (
	"automation/handlers"
	"log"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/change-product", handlers.ChangeProductHandler)
	http.HandleFunc("/delete-product", handlers.DeleteProductHandler)
	http.HandleFunc("/add-worker", handlers.AddWorkerHandler)
	http.HandleFunc("/get-orders", handlers.GetOrders)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/add-product", handlers.AddProductHandler)
	http.HandleFunc("/purchase-history", handlers.PurchaseHistoryHandler)
	http.HandleFunc("/total-purchases", handlers.TotalPurchasesHandler)
	http.HandleFunc("/get-users", handlers.GetUsersHandler)
	http.HandleFunc("/get-all-products", handlers.GetAllProductsHandler)   
	http.HandleFunc("/get-product-image/", handlers.GetProductImageHandler) 
	http.HandleFunc("/get-clients", handlers.GetClientsHandler)
	http.HandleFunc("/get-product-by-name/product", handlers.GetProductByName)
	http.HandleFunc("/get-order-by-name/order", handlers.GetOrdersByClientName)
	http.HandleFunc("/get-client-by-name/client", handlers.GetClientsByName)
	http.HandleFunc("/get-worker-by-name/worker", handlers.GetWorkersByName)
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(http.DefaultServeMux)))
}
