package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"automation/db"
)

func AddProductHandler(res http.ResponseWriter, req *http.Request) {
	var (
		status int
		err    error
	)

	defer func() {
		if err != nil {
			http.Error(res, err.Error(), status)
		}
	}()

	if err = req.ParseMultipartForm(32 << 20); err != nil {
		status = http.StatusInternalServerError
		return
	}
	fmt.Println("Form parsed successfully")

	name := req.FormValue("name")
	priceStr := req.FormValue("price")

	if name == "" {
		err = fmt.Errorf("Product name is required")
		status = http.StatusBadRequest
		return
	}

	var price float64
	if priceStr != "" {
		price, err = strconv.ParseFloat(priceStr, 64)
		if err != nil {
			err = fmt.Errorf("Invalid price format")
			status = http.StatusBadRequest
			return
		}
	} else {
		err = fmt.Errorf("Price is required")
		status = http.StatusBadRequest
		return
	}

	file, fileHeader, err := req.FormFile("photo")
	if err != nil {
		err = fmt.Errorf("Error retrieving file: %v", err)
		status = http.StatusInternalServerError
		return
	}
	defer file.Close()

	uploadDir := "./uploads/"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	filePath := uploadDir + fileHeader.Filename
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

	dbConn, err := db.ConnectDB()
	if err != nil {
		err = fmt.Errorf("Database connection error: %v", err)
		status = http.StatusInternalServerError
		return
	}
	defer dbConn.Close()

	_, err = dbConn.Exec("INSERT INTO product (name, price, photo) VALUES (?, ?, ?)", name, price, filePath)
	if err != nil {
		err = fmt.Errorf("Error inserting product into database: %v", err)
		status = http.StatusInternalServerError
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Product added successfully"))
}
