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
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	database, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer database.Close()

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		log.Println("Error parsing form data:", err)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		log.Println("Product ID is missing")
		return
	}

	newName := r.FormValue("new_name")
	newCategory := r.FormValue("new_category")
	newBrand := r.FormValue("new_brand")
	newStock := r.FormValue("new_stock")
	newPrice := r.FormValue("new_price")

	var currentName, currentPhoto string
	query := `SELECT name, photo FROM product WHERE id = ?`
	err = database.QueryRow(query, id).Scan(&currentName, &currentPhoto)
	if err != nil {
		http.Error(w, "Failed to retrieve product details", http.StatusInternalServerError)
		log.Println("Error retrieving product details:", err)
		return
	}

	var newPhotoPath string
	var photoNameDB string
	file, fileHeader, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()

		extension := getFileExtension(fileHeader.Filename)
		fileName := fmt.Sprintf("%s-%s%s", sanitizeFileName(newName), id, extension)
		photoNameDB = fmt.Sprintf("%s-%s%s", sanitizeFileName(newName), id, extension)
		newPhotoPath = fmt.Sprintf("./uploads/%s", fileName)

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

		if currentPhoto != "" && currentPhoto != newPhotoPath {

			err = os.Remove(fmt.Sprintf("./uploads/%s", currentPhoto))
			if err != nil {
				log.Println("Error deleting old photo:", err)
			}
		}
	} else if newName != "" && currentPhoto != "" {
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
		newPhotoPath = currentPhoto

	}

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

func getFileExtension(filename string) string {
	return filepath.Ext(filename)
}
func sanitizeFileName(name string) string {
	re := regexp.MustCompile(`[^\w\-_]+`)
	return re.ReplaceAllString(name, "_")
}

func nullableString(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}
