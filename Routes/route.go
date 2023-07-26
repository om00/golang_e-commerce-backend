package Routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/om00/golang-ecommerce/Controllers"
)

func UserRoutes(router *mux.Router, db *sql.DB) {

	router.HandleFunc("/users/sign-up", func(w http.ResponseWriter, r *http.Request) {
		Controllers.Signup(w, r, db)
	}).Methods("POST")
	router.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		Controllers.Login(w, r, db)
	}).Methods("POST")
	router.HandleFunc("/admin/addproduct", func(w http.ResponseWriter, r *http.Request) {
		Controllers.ProductViewAdmin(w, r, db)
	}).Methods("POST")
	router.HandleFunc("/users/prdouctview", func(w http.ResponseWriter, r *http.Request) {
		Controllers.SearchProduct(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/users/search", func(w http.ResponseWriter, r *http.Request) {
		Controllers.SearchProductByQuery(w, r, db)
	}).Methods("GET")
}
