// package handlers

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"automation/db"
// )

// func AddProductHandler(res http.ResponseWriter, req *http.Request) {
// 	var (
// 		status int
// 		err    error
// 	)

// 	// Убедимся, что запрос это POST
// 	if req.Method != http.MethodPost {
// 		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	defer func() {
// 		if err != nil {
// 			http.Error(res, err.Error(), status)
// 		}
// 	}()

// 	// Разбор данных из формы (включая файл)
// 	if err = req.ParseMultipartForm(32 << 20); err != nil {
// 		status = http.StatusInternalServerError
// 		return
// 	}
// 	fmt.Println("Form parsed successfully")

// 	// Извлечение данных формы
// 	name := req.FormValue("name")
// 	priceStr := req.FormValue("price")
// 	discript := req.FormValue("discript")
// 	idCategoriesStr := req.FormValue("idCategories")
// 	idBrandsStr := req.FormValue("idBrands")
// 	qualityStr := req.FormValue("quality")

// 	// Печать полученных данных для отладки
// 	fmt.Println("Name:", name)
// 	fmt.Println("Price:", priceStr)
// 	fmt.Println("Description:", discript)
// 	fmt.Println("Category ID:", idCategoriesStr)
// 	fmt.Println("Brand ID:", idBrandsStr)
// 	fmt.Println("Quality:", qualityStr)

// 	// Проверка обязательных полей
// 	if name == "" {
// 		err = fmt.Errorf("Product name is required")
// 		status = http.StatusBadRequest
// 		return
// 	}
// 	if discript == "" {
// 		err = fmt.Errorf("Product description is required")
// 		status = http.StatusBadRequest
// 		return
// 	}
// 	if priceStr == "" {
// 		err = fmt.Errorf("Price is required")
// 		status = http.StatusBadRequest
// 		return
// 	}

// 	// Конвертация price
// 	price, convErr := strconv.Atoi(priceStr)
// 	if convErr != nil {
// 		err = fmt.Errorf("Invalid price format")
// 		status = http.StatusBadRequest
// 		return
// 	}

// 	// Обработка файла
// 	file, fileHeader, fileErr := req.FormFile("photo")
// 	if fileErr != nil {
// 		err = fmt.Errorf("Error retrieving file: %v", fileErr)
// 		status = http.StatusBadRequest
// 		return
// 	}
// 	defer file.Close()

// 	// Путь для загрузки файла
// 	uploadDir := "./uploads/"
// 	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
// 		err := os.Mkdir(uploadDir, os.ModePerm)
// 		if err != nil {
// 			status = http.StatusInternalServerError
// 			return
// 		}
// 	}

// 	// Сохраняем файл на диск
// 	filePath := uploadDir + fileHeader.Filename
// 	outfile, err := os.Create(filePath)
// 	if err != nil {
// 		err = fmt.Errorf("Error saving file: %v", err)
// 		status = http.StatusInternalServerError
// 		return
// 	}
// 	defer outfile.Close()

// 	_, err = io.Copy(outfile, file)
// 	if err != nil {
// 		err = fmt.Errorf("Error copying file: %v", err)
// 		status = http.StatusInternalServerError
// 		return
// 	}
// 	fmt.Printf("File successfully uploaded: %s\n", filePath)

// 	// Конвертация дополнительных полей, если они присутствуют
// 	var idCategories *int
// 	var idBrands *int
// 	var quality *int

// 	if idCategoriesStr != "" {
// 		val, convErr := strconv.Atoi(idCategoriesStr)
// 		if convErr == nil {
// 			idCategories = &val
// 		}
// 	}

// 	if idBrandsStr != "" {
// 		val, convErr := strconv.Atoi(idBrandsStr)
// 		if convErr == nil {
// 			idBrands = &val
// 		}
// 	}

// 	if qualityStr != "" {
// 		val, convErr := strconv.Atoi(qualityStr)
// 		if convErr == nil {
// 			quality = &val
// 		}
// 	}

// 	// Подключение к базе данных
// 	dbConn, err := db.ConnectDB()
// 	if err != nil {
// 		err = fmt.Errorf("Database connection error: %v", err)
// 		status = http.StatusInternalServerError
// 		return
// 	}
// 	defer dbConn.Close()

// 	// SQL-запрос на добавление продукта в базу данных
// 	query := `INSERT INTO product (name, photo, discript, idCategories, idBrands, quality, price)
// 	          VALUES (?, ?, ?, ?, ?, ?, ?)`

// 	_, err = dbConn.Exec(query, name, filePath, discript, idCategories, idBrands, quality, price)
// 	if err != nil {
// 		err = fmt.Errorf("Error inserting product into database: %v", err)
// 		status = http.StatusInternalServerError
// 		return
// 	}

// 	// Ответ клиенту о успешном добавлении продукта
// 	res.WriteHeader(http.StatusOK)
// 	res.Write([]byte("Product added successfully"))
// }
package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"automation/db"
)

func AddProductHandler(res http.ResponseWriter, req *http.Request) {
	var (
		status int
		err    error
	)

	// Убедимся, что запрос это POST
	if req.Method != http.MethodPost {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	defer func() {
		if err != nil {
			http.Error(res, err.Error(), status)
		}
	}()

	// Разбор данных из формы (включая файл)
	if err = req.ParseMultipartForm(32 << 20); err != nil {
		status = http.StatusInternalServerError
		return
	}
	fmt.Println("Form parsed successfully")

	// Извлечение данных формы
	name := req.FormValue("name")
	priceStr := req.FormValue("price")
	discript := req.FormValue("discript")
	idCategoriesStr := req.FormValue("idCategories")
	idBrandsStr := req.FormValue("idBrands")
	qualityStr := req.FormValue("quality")

	// Печать полученных данных для отладки
	fmt.Println("Name:", name)
	fmt.Println("Price:", priceStr)
	fmt.Println("Description:", discript)
	fmt.Println("Category ID:", idCategoriesStr)
	fmt.Println("Brand ID:", idBrandsStr)
	fmt.Println("Quality:", qualityStr)

	// Проверка обязательных полей
	if name == "" {
		err = fmt.Errorf("Product name is required")
		status = http.StatusBadRequest
		return
	}
	if discript == "" {
		err = fmt.Errorf("Product description is required")
		status = http.StatusBadRequest
		return
	}
	if priceStr == "" {
		err = fmt.Errorf("Price is required")
		status = http.StatusBadRequest
		return
	}

	// Конвертация price
	price, convErr := strconv.Atoi(priceStr)
	if convErr != nil {
		err = fmt.Errorf("Invalid price format")
		status = http.StatusBadRequest
		return
	}

	// Обработка файла
	file, fileHeader, fileErr := req.FormFile("photo")
	if fileErr != nil {
		err = fmt.Errorf("Error retrieving file: %v", fileErr)
		status = http.StatusBadRequest
		return
	}
	defer file.Close()

	// Формируем новое имя файла
	now := time.Now()
	date := now.Format("20060102")
	timeFormatted := now.Format("150405")
	idCategories, _ := strconv.Atoi(idCategoriesStr)
	idBrands, _ := strconv.Atoi(idBrandsStr)
	ext := fileHeader.Filename[len(fileHeader.Filename)-4:] // Получаем расширение файла

	// Формируем путь для загрузки файла
	uploadDir := "./uploads/"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			status = http.StatusInternalServerError
			return
		}
	}

	// Формируем имя файла
	newFileName := fmt.Sprintf("/%s_%s_%s_%d_%d%s", name, date, timeFormatted, idCategories, idBrands, ext)
	filePath := uploadDir + newFileName

	// Сохраняем файл на диск
	outfile, err := os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("Error saving file: %v", err)
		status = http.StatusInternalServerError
		return
	}
	defer outfile.Close()

	_, err = io.Copy(outfile, file)
	if err != nil {
		err = fmt.Errorf("Error copying file: %v", err)
		status = http.StatusInternalServerError
		return
	}
	fmt.Printf("File successfully uploaded: %s\n", filePath)

	// Конвертация дополнительных полей, если они присутствуют
	var idCategoriesInt *int
	var idBrandsInt *int
	var qualityInt *int

	if idCategoriesStr != "" {
		val, convErr := strconv.Atoi(idCategoriesStr)
		if convErr == nil {
			idCategoriesInt = &val
		}
	}

	if idBrandsStr != "" {
		val, convErr := strconv.Atoi(idBrandsStr)
		if convErr == nil {
			idBrandsInt = &val
		}
	}

	if qualityStr != "" {
		val, convErr := strconv.Atoi(qualityStr)
		if convErr == nil {
			qualityInt = &val
		}
	}

	// Подключение к базе данных
	dbConn, err := db.ConnectDB()
	if err != nil {
		err = fmt.Errorf("Database connection error: %v", err)
		status = http.StatusInternalServerError
		return
	}
	defer dbConn.Close()

	// SQL-запрос на добавление продукта в базу данных
	query := `INSERT INTO product (name, photo, discript, idCategories, idBrands, quality, price)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = dbConn.Exec(query, name, newFileName, discript, idCategoriesInt, idBrandsInt, qualityInt, price)
	if err != nil {
		err = fmt.Errorf("Error inserting product into database: %v", err)
		status = http.StatusInternalServerError
		return
	}

	// Ответ клиенту о успешном добавлении продукта
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Product added successfully"))
}
