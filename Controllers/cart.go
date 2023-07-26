package Controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Application struct {
	db *sql.DB
}

func NewApplication(db *sqlDB) *Application {
	return &Application{
		db: db,
	}
}

func (app *Application) AddToCart(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()
	user_id := query_values.Get("userId")
	if user_id == "" {
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return

	}

	product_id := query_values.Get("productId")

	if product_id == "" {
		http.Error(w, "Product is missing", http.StatusBadRequest)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := Database.AddProductToCart(ctx, app, user_id, proudct_id)

	if err != nil {
		http.Error(w, err, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print("Successfully addeded to the cart")
}

func (app *Application) RemoveItem(w http.ResponseWriter, r *http.Request) {

	query_values := r.URL.Query()
	user_id := query_values.Get("userId")
	if user_id == "" {
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return

	}

	product_id := query_values.Get("productId")

	if product_id == "" {
		http.Error(w, "Product is missing", http.StatusBadRequest)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := Database.RemoveItemFromCart(ctx, app, user_id, product_id)

	if err != nil {
		http.Error(w, err, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print("Successfully item removed form the cart")

}

func (app *Application) GetItemFromCart(w http.ResponseWriter, r *http.Request) {
	filer_values := r.URL.Query()
	user_id = filer_values.Get("id")

	ctx = r.context()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var user Model.User
	query := "Select id,userCart from User where id=?"
	row, err := DB.QueryRow(query, user_id)

	if err != nil {
		log.Println("error in execution of the query")
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	err := row.Scan(&user.ID, &user.User_Cart)

	if err != nil {
		log.Println("Error in reading the data", err)
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.Header.Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(user.User_Cart)

	if err != nil {
		log.Println("error while encoding in json ", err)
		http.Error(w, err, http.StatusInternalServerError)
	}

}

func (app *Application) BuyFromCart(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()
	user_id := query_values.Get("userId")
	if user_id == "" {
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return

	}

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := Database.BuyIteamFromCart(ctx, app, user_id)
	if err != nil {
		http.Error(w, err, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print("Success")

}

func (app *Application) InstantBuy(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()
	user_id := query_values.Get("userId")
	if user_id == "" {
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return

	}

	product_id := query_values.Get("productId")

	if product_id == "" {
		http.Error(w, "Product is missing", http.StatusBadRequest)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := Database.InstantBuyer(ctx, app, user_id, product_id)
	if err != nil {
		http.Error(w, err, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print("success")
}
