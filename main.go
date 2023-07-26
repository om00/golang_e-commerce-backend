package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/om00/golang-ecommerce/Controllers"
	"github.com/om00/golang-ecommerce/Database"
	"github.com/om00/golang-ecommerce/Middleware"
	"github.com/om00/golang-ecommerce/Routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	db, err := Database.InitDB()
	if err != nil {
		log.Println("error while when connecting ot the database", err)
	}
	app := Controllers.NewApplication(db)

	r := mux.NewRouter()

	Routes.UserRoutes(r, db)
	r.Use(Middleware.Authentication)

	r.HandleFunc("/addAddress", func(w http.ResponseWriter, r *http.Request) {
		Controllers.AddAddress(w, r, db)
	}).Methods("POST")
	r.HandleFunc("/DeleteAddress", func(w http.ResponseWriter, r *http.Request) {
		Controllers.DeleteAddress(w, r, db)
	}).Methods("GET")
	r.HandleFunc("/EditAddress", func(w http.ResponseWriter, r *http.Request) {
		Controllers.EditAddress(w, r, db)
	}).Methods("PUT")
	r.HandleFunc("/ListCart", func(w http.ResponseWriter, r *http.Request) {
		Controllers.GetItemFromCart(w, r, db)
	}).Methods("GET")
	r.HandleFunc("/addtocart", app.AddToCart).Methods("GET")
	r.HandleFunc("/removeitem", app.RemoveItem).Methods("GET")
	r.HandleFunc("/cartcheckout", app.BuyFromCart).Methods("GET")
	r.HandleFunc("/instantbuy", app.InstantBuy).Methods("GET")

	fmt.Println("Service is running on port no :", port)
	http.ListenAndServe(":9090", r)

}
