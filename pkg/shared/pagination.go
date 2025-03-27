package shared

import (
	"fmt"

	"github.com/heshanu/go-service/pkg/models"
)

type ProductResponse struct {
	TotalProductCount int
	Page              int
	Limit             int
	ProductList       []models.Product
}

// keep in mind renaming function Public method -start with UpperCase
// private method-lowerCase
func ProductPagination(productList []models.Product, page int, limit int) (ProductResponse, error) {
	if page <= 0 || limit <= 0 {
		return ProductResponse{}, fmt.Errorf("page and limit must be greater than zero")
	}

	// Calculate pagination offsets
	start := (page - 1) * limit
	end := start + limit

	// Ensure end does not exceed the length of the books slice
	if end > len(productList) {
		end = len(productList)
	}

	// Get the paginated subset of books
	paginatedProduct := productList[start:end]

	// Create the paginated response
	response := ProductResponse{
		TotalProductCount: len(productList),
		ProductList:       paginatedProduct,
		Page:              page,
		Limit:             limit,
	}
	return response, nil
}
