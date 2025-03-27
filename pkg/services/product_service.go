package services

import (
	"fmt"

	"github.com/heshanu/go-service/pkg/config"
	"github.com/heshanu/go-service/pkg/models"
	"github.com/heshanu/go-service/pkg/shared"
	"gorm.io/gorm"
)

func GetAllProduct() ([]models.Product, error) {
	var productList []models.Product
	db := config.GetDB()
	if err := db.Find(&productList).Error; err != nil {
		return nil, err
	}
	return productList, nil
}

// CreateBook creates a new book record in the database
func CreateProduct(book *models.Product) error {
	db := config.GetDB()
	if db == nil {
		return gorm.ErrInvalidTransaction // or another appropriate error
	}

	if err := db.Create(book).Error; err != nil {
		return err
	}
	return nil
}

func GetBookProductById(id int64) (*models.Product, error) {
	var product models.Product

	db := config.GetDB()
	if err := db.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // or return a custom error if you prefer
		}
		return nil, err
	}
	return &product, nil
}

func DeleteProductById(id int64) error {
	db := config.GetDB()
	if err := db.Delete(&models.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProductById(id int64, updatedFields *models.Product) error {
	db := config.GetDB()

	// Find the book by ID
	var product models.Product
	if err := db.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	// Check for non-zero fields in updatedFields and update them
	if updatedFields.Name != "" {
		product.Name = updatedFields.Name
	}

	if updatedFields.Description != "" {
		product.Description = updatedFields.Description
	}

	if updatedFields.Price != 0 {
		product.Price = updatedFields.Price
	}

	if updatedFields.Quantity != 0 {
		product.Quantity = updatedFields.Quantity
	}

	// Save the updated book to the database
	if err := db.Save(&product).Error; err != nil {
		return err
	}

	return nil
}

func SearchByKeyWord(productList []models.Product, keyWord string) ([]models.Product, error) {
	if keyWord == "" {
		fmt.Printf("Keyword Cannot null")
		return []models.Product{}, fmt.Errorf("keyword cannoit be null")
	}
	db := config.GetDB()
	if err := db.Find(&productList).Error; err != nil {
		return []models.Product{}, err
	}

	searchResponse := shared.SearchBooksConcurrently(productList, keyWord)
	return searchResponse, nil

}
