package handlers

import (
	"automation/db"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func ChangeProductHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Подключаемся к базе данных
	database, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer database.Close()

	// Парсим мультипарт-запрос
	err = r.ParseMultipartForm(10 << 20) // Ограничение на размер файла (10 МБ)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		log.Println("Error parsing form data:", err)
		return
	}

	// Получаем ID продукта
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		log.Println("Product ID is missing")
		return
	}

	// Получаем остальные данные
	newName := r.FormValue("new_name")
	newCategory := r.FormValue("new_category")
	newBrand := r.FormValue("new_brand")
	newStock := r.FormValue("new_stock")
	newPrice := r.FormValue("new_price")

	// Получаем текущее имя продукта и путь к фото
	var currentName, currentPhoto string
	query := `SELECT name, photo FROM product WHERE id = ?`
	err = database.QueryRow(query, id).Scan(&currentName, &currentPhoto)
	if err != nil {
		http.Error(w, "Failed to retrieve product details", http.StatusInternalServerError)
		log.Println("Error retrieving product details:", err)
		return
	}

	// Работа с фото
	var newPhotoPath string
	var photoNameDB string
	file, fileHeader, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()

		// Формируем имя файла: новое_имя-id.формат
		extension := getFileExtension(fileHeader.Filename)
		fileName := fmt.Sprintf("%s-%s%s", sanitizeFileName(newName), id, extension)
		photoNameDB = fmt.Sprintf("%s-%s%s", sanitizeFileName(newName), id, extension)
		newPhotoPath = fmt.Sprintf("./uploads/%s", fileName)

		// Сохраняем новый файл
		out, err := os.Create(newPhotoPath)
		if err != nil {
			http.Error(w, "Failed to save the file", http.StatusInternalServerError)
			log.Println("Error saving file:", err)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Failed to write the file", http.StatusInternalServerError)
			log.Println("Error writing file:", err)
			return
		}

		// Удаляем старое фото, если оно существует
		if currentPhoto != "" && currentPhoto != newPhotoPath {

			err = os.Remove(fmt.Sprintf("./uploads/%s", currentPhoto))
			if err != nil {
				log.Println("Error deleting old photo:", err)
			}
		}
	} else if newName != "" && currentPhoto != "" {
		// Переименование фото, если изменилось имя
		oldPhotoPath := fmt.Sprintf("./uploads/%s", currentPhoto)
		extension := getFileExtension(oldPhotoPath)
		photoNameDB = fmt.Sprintf("/%s-%s%s", sanitizeFileName(newName), id, extension)
		newPhotoPath = fmt.Sprintf("./uploads/%s-%s%s", sanitizeFileName(newName), id, extension)

		err = os.Rename(oldPhotoPath, newPhotoPath)
		if err != nil {
			http.Error(w, "Failed to rename the photo", http.StatusInternalServerError)
			log.Println("Error renaming photo:", err)
			return
		}
	} else {
		// Если фото не менялось
		newPhotoPath = currentPhoto

	}

	// Формируем запрос к базе данных
	query = `
		UPDATE product 
		SET name = COALESCE(?, name),
			photo = COALESCE(?, photo),
			idCategories = COALESCE(?, idCategories),
			idBrands = COALESCE(?, idBrands),
			quality = COALESCE(?, quality),
			price = COALESCE(?, price)
		WHERE id = ?`
	_, err = database.Exec(query, nullableString(newName), nullableString(photoNameDB), nullableString(newCategory),
		nullableString(newBrand), nullableString(newStock), nullableString(newPrice), id)

	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		log.Println("Error updating product:", err)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

// Получение расширения файла
func getFileExtension(filename string) string {
	return filepath.Ext(filename)
}

// Очистка имени файла от недопустимых символов
func sanitizeFileName(name string) string {
	re := regexp.MustCompile(`[^\w\-_]+`)
	return re.ReplaceAllString(name, "_")
}

// Преобразование пустых строк в NULL
func nullableString(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}
