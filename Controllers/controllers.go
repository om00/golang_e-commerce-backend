package Controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/om00/golang-ecommerce/Database"
	"github.com/om00/golang-ecommerce/Models"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)

}

func verifyPassword(user_pass, given_pass string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(given_pass), []byte(user_pass))
	valid := true
	msg := ""

	if err != nil {
		valid = false
		msg = "Incorrect Password"

	}

	return valid, msg
}

// func (ctrl *Controller) ProductViewAdmin(w http.ResponseWriter, r *http.Request) {

// }

// func (ctrl *Controller) SearchProduct(w http.ResponseWriter, r *http.Request) {

// 	var ProductList []Models.Product

// 	query := "Select id,productName,price,rating,image from Product"
// 	rows, err := db.Query(query)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	for rows.Next() {
// 		var row Models.Product
// 		err := rows.Scan(&row.ID, &row.Product_Name, &row.Price, &row.Rating, &row.Image)
// 		if err != nil {
// 			fmt.Println("Error while scaning menthod")
// 			return
// 		}

// 		ProductList = append(ProductList, row)

// 	}

// 	w.Header().Set("Content-Type", "application/json")

// 	err = json.NewEncoder(w).Encode(ProductList)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// }
func (ctrl *Controller) SearchProductByQuery(w http.ResponseWriter, r *http.Request) {

	req_values := r.URL.Query()

	var ctx, cancel = context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	name_filter := req_values.Get("name")
	fitler := Models.ProductQuery{
		Name: name_filter,
	}

	product, err := ctrl.db.GetProuducts(ctx, fitler)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(product)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

type Controller struct {
	db *Database.DB
	//can dd more
}

func NewController(db *Database.DB) *Controller {
	return &Controller{
		db: db,
	}
}
