package Controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (ctrl *Controller) AddToCart(w http.ResponseWriter, r *http.Request) {
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

	err = ctrl.db.AddProductToCart(ctx, user_id, product_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print("Successfully addeded to the cart")
}

func (ctrl *Controller) RemoveItem(w http.ResponseWriter, r *http.Request) {

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

	err = ctrl.db.RemoveItemFromCart(ctx, user_id, product_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print("Successfully item removed form the cart")

}

// func (ctrl *Controller) GetItemFromCart(w http.ResponseWriter, r *http.Request) {
// 	filer_values := r.URL.Query()
// 	user_id := filer_values.Get("id")

// 	var user Models.User
// 	query := "Select id,userCart from User where id=?"
// 	err := ctrl.d.QueryRow(query, user_id).Scan(&user.ID, &user.UserCart)

// 	if err != nil {
// 		log.Println("error in execution of the query")
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if err == sql.ErrNoRows {
// 		log.Println("No data is present with this id")
// 		http.Error(w, err.Error(), http.StatusOK)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	err = json.NewEncoder(w).Encode(user.UserCart)

// 	if err != nil {
// 		log.Println("error while encoding in json ", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}

// }

func (ctrl *Controller) BuyFromCart(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()
	user_id, err := strconv.ParseInt(query_values.Get("userId"), 10, 54)
	if user_id == 0 || err != nil {
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return

	}

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = ctrl.db.BuyIteamFromCart(ctx, user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print("Success")

}

func (ctrl *Controller) InstantBuy(w http.ResponseWriter, r *http.Request) {
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

	err = ctrl.db.InstantBuyer(ctx, product_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print("success")
}
