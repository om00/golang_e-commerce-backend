package Controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/om00/golang-ecommerce/Database"
	"github.com/om00/golang-ecommerce/Models"
)

type Application struct {
	db *sql.DB
}

func NewApplication(db *sql.DB) *Application {
	return &Application{
		db: db,
	}
}

func (app *Application) AddToCart(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()
	user_id, err := strconv.ParseInt(query_values.Get("userId"), 10, 64)
	if user_id == 0 || err != nil {
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return

	}

	product_id, err := strconv.ParseInt(query_values.Get("productId"), 10, 64)

	if product_id == 0 || err != nil {
		http.Error(w, "Product is missing", http.StatusBadRequest)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Database.AddProductToCart(ctx, app.db, user_id, product_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print("Successfully addeded to the cart")
}

func (app *Application) RemoveItem(w http.ResponseWriter, r *http.Request) {

	query_values := r.URL.Query()
	user_id, err := strconv.ParseInt(query_values.Get("userId"), 10, 64)
	if user_id == 0 || err != nil {
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return

	}

	product_id, err := strconv.ParseInt(query_values.Get("productId"), 10, 64)

	if product_id == 0 || err != nil {
		http.Error(w, "Product is missing", http.StatusBadRequest)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Database.RemoveItemFromCart(ctx, app.db, user_id, product_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print("Successfully item removed form the cart")

}

func GetItemFromCart(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	filer_values := r.URL.Query()
	user_id := filer_values.Get("id")

	var user Models.User
	query := "Select id,userCart from User where id=?"
	err := db.QueryRow(query, user_id).Scan(&user.ID, &user.UserCart)

	if err != nil {
		log.Println("error in execution of the query")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err == sql.ErrNoRows {
		log.Println("No data is present with this id")
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user.UserCart)

	if err != nil {
		log.Println("error while encoding in json ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (app *Application) BuyFromCart(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()
	user_id, err := strconv.ParseInt(query_values.Get("userId"), 10, 54)
	if user_id == 0 || err != nil {
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return

	}

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Database.BuyIteamFromCart(ctx, app.db, user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print("Success")

}

func (app *Application) InstantBuy(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()
	user_id, err := strconv.ParseInt(query_values.Get("userId"), 10, 64)
	if user_id == 0 || err != nil {
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return

	}

	product_id, err := strconv.ParseInt(query_values.Get("productId"), 10, 64)

	if product_id == 0 || err != nil {
		http.Error(w, "Product is missing", http.StatusBadRequest)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Database.InstantBuyer(ctx, app.db, product_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print("success")
}
