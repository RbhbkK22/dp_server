package main

import (
	"automation/handlers"
	"log"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем доступ с любого домена
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Разрешаем необходимые методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Разрешаем заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Если запрос типа OPTIONS (pre-flight), просто возвращаем статус 200
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Передаем запрос дальше
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
	http.HandleFunc("/get-all-products", handlers.GetAllProductsHandler)    // Обработчик для получения товаров
	http.HandleFunc("/get-product-image/", handlers.GetProductImageHandler) // Обработчик для получения изображения
	http.HandleFunc("/get-clients", handlers.GetClientsHandler)
	http.HandleFunc("/get-product-by-name/product", handlers.GetProductByName)
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(http.DefaultServeMux)))
}
