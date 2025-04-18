package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/om00/golang-ecommerce/Controllers"
	"github.com/om00/golang-ecommerce/Database"
	"github.com/om00/golang-ecommerce/Routes"
	"github.com/om00/golang-ecommerce/Token"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Token.SECRET_KEY = os.Getenv("SECRET_KEY")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	dbObj, err := Database.NewDB()
	if err != nil {
		log.Fatal("error while when connecting ot the database", err)
	}

	err = dbObj.RunMigrations()
	if err != nil {
		log.Fatal("error while running migrations", err)
	}

	ctrl := Controllers.NewController(dbObj)

	r := mux.NewRouter()

	api := Routes.NewApi(ctrl, r)
	api.UserRoutes()
	api.RoutesWithMiddleWare()

	fmt.Println("Service is running on port no :", port)
	http.ListenAndServe(":9090", r)

}
