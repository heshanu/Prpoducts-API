package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/heshanu/go-service/pkg/models"
	"github.com/heshanu/go-service/pkg/services"
	"github.com/heshanu/go-service/pkg/shared"
	"github.com/heshanu/go-service/pkg/utils"
	"gorm.io/gorm"
)

var (
	NewBook models.Product
	//selectedBook services.Book
)

func GetAllProducts(w http.ResponseWriter, request *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(request.URL.Query().Get("limit"))

	productList, err := services.GetAllProduct()

	if err != nil {
		http.Error(w, "Cannot get all products", http.StatusFound)
	}

	var resProduct shared.ProductResponse

	resProduct, err = shared.ProductPagination(productList, page, limit)

	if err != nil {
		http.Error(w, "Issue with pagination process", http.StatusConflict)
	}

	res, _ := json.Marshal(resProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["productId"]
	ID, err := strconv.ParseInt(productId, 0, 0)

	if err != nil {
		http.Error(w, "Parse null  product Id", http.StatusBadRequest)
	}

	productDetail, err := services.GetBookProductById(ID)

	if err != nil {
		http.Error(w, "Cannot find product by this value", http.StatusBadRequest)
	}

	res, _ := json.Marshal(productDetail)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into the createBook struct
	var createProduct models.Product
	if err := utils.ParseBody(r, &createProduct); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if createProduct.Name == "" || createProduct.Description == "" ||
		createProduct.Price == 0 || createProduct.Quantity == 0 {
		http.Error(w, "Name, description, and price,quantity are required fields", http.StatusBadRequest)
		return
	}

	// Create the book in the database
	if err := services.CreateProduct(&createProduct); err != nil {
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		return
	}

	// Marshal the created book into JSON
	res, err := json.Marshal(createProduct)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse the book ID from the URL
	vars := mux.Vars(r)
	bookId := vars["productId"]
	ID, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	// Parse the request body to get the updated book details
	var updatedBook models.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the book in the database
	if err := services.UpdateProductById(ID, &updatedBook); err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update book", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the updated book details
	res, err := json.Marshal(updatedBook)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["productId"]
	ID, err := strconv.ParseInt(productId, 10, 64)

	if err != nil {
		http.Error(w, "Parse null productId", http.StatusBadRequest)
	}

	er := services.DeleteProductById(ID)
	if er != nil {
		http.Error(w, "Error while deleting book", http.StatusBadRequest)
	}

	res, _ := json.Marshal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

/*
func SearchBookByKeyword(w http.ResponseWriter, request *http.Request) {
	q := request.URL.Query().Get("q")

	productList, err := services.GetAllProduct()

	if err != nil {
		http.Error(w, "Cannot get all products", http.StatusFound)
	}

	var result []models.Product
	result, err = services.SearchByKeyWord(productList, q)

	if err != nil {
		http.Error(w, "Cannot get keyword related products", http.StatusFound)
	}

	res, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}*/

func SearchBookByKeyword(w http.ResponseWriter, request *http.Request) {
	q := request.URL.Query().Get("q")

	productList, err := services.GetAllProduct()
	if err != nil {
		log.Println("Error fetching products:", err)
		http.Error(w, "Cannot get all products", http.StatusInternalServerError)
		return
	}

	result := shared.SearchBooksConcurrently(productList, q)
	// if err != nil {
	// 	log.Println("Error searching products by keyword:", err)
	// 	http.Error(w, "Cannot get keyword related products", http.StatusInternalServerError)
	// 	return
	// }

	res, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
