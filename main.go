package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/om00/golang-ecommerce/Controllers"
	"github.com/om00/golang-ecommerce/Database"
	"github.com/om00/golang-ecommerce/Routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	dbObj, err := Database.NewDB()
	if err != nil {
		log.Println("error while when connecting ot the database", err)
		os.Exit(0)
	}

	ctrl := Controllers.NewController(dbObj)

	r := mux.NewRouter()

	api := Routes.NewApi(ctrl, r)
	api.UserRoutes()
	api.RoutesWithMiddleWare()

	fmt.Println("Service is running on port no :", port)
	http.ListenAndServe(":9090", r)

}
