package main

import (
	"os"

	"github.com/gorilla/mux"
	"github.com/om00/golang-ecommerce/Controllers"
	"github.com/om00/golang-ecommerce/Middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	db, err := Database.initDB()
	app := Controllers.NewApplication(db)

	r := mux.NewRouter()

	Routes.userRoutes(r)
	r.Use(Middleware.Authentication())

	r.Handlefunc("/addtocart", app.AddToCart).Methods("GET")
	r.Handlefunc("/removeitem", app.RemoveItem).Methods("GET")
	r.Handlefunc("/cartcheckout", app.BuyFromCart).Methods("GET")
	r.Handlefunc("/instantbuy", app.InstatnBuy).Methods("GET")

	log.fatal(r.Run(":", port))
}
