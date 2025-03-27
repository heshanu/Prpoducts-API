package routes

import (
	"github.com/gorilla/mux"
	"github.com/heshanu/go-service/pkg/controllers"
)

var RegisterBookStoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/product", controllers.CreateProduct).Methods("POST")
	router.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")

	router.HandleFunc("/product/{productId}", controllers.GetProductById).Methods("GET")

	router.HandleFunc("/products/productSearch", controllers.SearchBookByKeyword).Methods("GET")

	router.HandleFunc("/product/{productId}", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/product/{productId}", controllers.DeleteProductById).Methods("DELETE")
}
