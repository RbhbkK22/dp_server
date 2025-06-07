package handlers

import (
	"automation/db"
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/disintegration/imaging"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
)

type Product struct {
	Category string
	Brand    string
	Name     string
	Photo    string
	Price    float64
}

func prepareImageForExcel(originalPath string) (string, error) {
	// Открываем файл
	file, err := os.Open(originalPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Пытаемся декодировать
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return "", fmt.Errorf("не удалось определить формат: %v", err)
	}

	// Поддерживаемые форматы Excel
	if format == "jpeg" || format == "png" || format == "gif" {
		return originalPath, nil // Всё ок
	}

	// Нужно переконвертировать
	img, err := imaging.Open(originalPath)
	if err != nil {
		return "", fmt.Errorf("ошибка открытия изображения через imaging: %v", err)
	}

	tempPng := strings.TrimSuffix(originalPath, filepath.Ext(originalPath)) + "_converted.png"
	err = imaging.Save(img, tempPng)
	if err != nil {
		return "", fmt.Errorf("ошибка сохранения PNG: %v", err)
	}

	return tempPng, nil
}

func ExporpHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer database.Close()

	rows, err := database.Query(`
		SELECT c.name as category, b.name as brand, p.name, p.photo, p.price
		FROM product p
		JOIN categories c ON p.idCategories = c.id
		JOIN brands b ON p.idBrands = b.id
		ORDER BY c.name, b.name
	`)

	if err != nil {
		http.Error(w, "Ошибка запроса", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var pr Product
		if err := rows.Scan(&pr.Category, &pr.Brand, &pr.Name, &pr.Photo, &pr.Price); err != nil {
			http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
			return
		}
		products = append(products, pr)
	}

	sort.Slice(products, func(i, j int) bool {
		if products[i].Category == products[j].Category {
			if products[i].Brand == products[j].Brand {
				return products[i].Name < products[j].Name
			}
			return products[i].Brand < products[j].Brand
		}
		return products[i].Category < products[j].Category

	})

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	row := 1
	currCat := ""
	currBrand := ""

	// Стиль для категории (темно синий фон, белый текст, жирный шрифт)
	categoryStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#3D6CE5"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
		},
	})

	// Стиль для бренда (ораньжевый фон, белый текст)
	brandStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#F07B1D"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
		},
	})

	// // Стиль для категории (синий фон, черный текст, жирный шрифт)
	// categoryStyle, _ := f.NewStyle(&excelize.Style{
	// 	Fill: excelize.Fill{
	// 		Type:    "pattern",
	// 		Color:   []string{"#DDEEFF"},
	// 		Pattern: 1,
	// 	},
	// 	Font: &excelize.Font{
	// 		Bold:  true,
	// 		Color: "#000000",
	// 	},
	// })

	// // Стиль для бренда (серый фон, черный текст)
	// brandStyle, _ := f.NewStyle(&excelize.Style{
	// 	Fill: excelize.Fill{
	// 		Type:    "pattern",
	// 		Color:   []string{"#EEEEEE"},
	// 		Pattern: 1,
	// 	},
	// 	Font: &excelize.Font{
	// 		Bold:  true,
	// 		Color: "#000000",
	// 	},
	// })

	for _, p := range products {
		if p.Category != currCat {
			currCat = p.Category
			currBrand = ""
			startCell := fmt.Sprintf("A%d", row)
			endCell := fmt.Sprintf("C%d", row)

			// Объединяем ячейки A-C в этой строке
			if err := f.MergeCell(sheet, startCell, endCell); err != nil {
				log.Println("Ошибка объединения ячеек категории:", err)
			}

			f.SetCellValue(sheet, startCell, currCat)
			f.SetCellStyle(sheet, startCell, endCell, categoryStyle)
			row++
		}
		if p.Brand != currBrand {
			currBrand = p.Brand
				startCell := fmt.Sprintf("A%d", row)
			endCell := fmt.Sprintf("C%d", row)

			// Объединяем ячейки A-C в этой строке
			if err := f.MergeCell(sheet, startCell, endCell); err != nil {
				log.Println("Ошибка объединения ячеек категории:", err)
			}
			f.SetCellValue(sheet, startCell, currBrand)
			f.SetCellStyle(sheet, startCell, endCell, brandStyle)
			row++

		}
		// Путь к изображению
		imgPath := "./uploads/" + p.Photo
		preparedImg, err := prepareImageForExcel(imgPath)
		if err == nil {
			f.SetColWidth(sheet, "A", "A", 20)
			f.SetRowHeight(sheet, row, 100)
			cell := fmt.Sprintf("A%d", row)
			err := f.AddPicture(sheet, cell, preparedImg, &excelize.GraphicOptions{
				AutoFit: true,
				OffsetX: 0, OffsetY: 0,
				ScaleX: 1, ScaleY: 1,
			})
			if err != nil {
				log.Println("Ошибка вставки картинки:", err)
			}
		} else {
			log.Println("Ошибка подготовки картинки:", err)
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Фото не вставлено")
		}
		// Название и цена
		f.SetColWidth(sheet, "B", "B", 40)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), p.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row) , fmt.Sprintf("%.2f руб.", p.Price))
		row++
	}

	// Папка экспорта
	if _, err := os.Stat("./exports"); os.IsNotExist(err) {
		_ = os.Mkdir("./exports", os.ModePerm)
	}

	filePath := "./exports/products.xlsx"
	if err := f.SaveAs(filePath); err != nil {
		http.Error(w, "Не удалось сохранить файл", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=products.xlsx")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, filePath)
}
