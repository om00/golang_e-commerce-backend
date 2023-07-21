package Routes

import (
	"github.com/gorilla/mux"
	"github.com/om00/golang-ecommerce/Controllers"
)

func userRoutes(router *mux.Router) {

	router.Handlefunc("/users/sign-up", Controllers.Signup).Methods("POST")
	router.Handlefunc("/users/login", Controllers.Login).Methods("POST")
	router.Handlefunc("/admin/addproduct", Controllers.ProductViewAdmin).Methods("POST")
	router.Handlefunc("/users/prdouctview", Controllers.SearchProduct).Methods("GET")
	router.Handlefunc("/users/search", Controllers.SearchProductByQuery).Methods("GET")
}
