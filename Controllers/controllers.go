package Controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/om00/golang-ecommerce/Database"
	"github.com/om00/golang-ecommerce/Models"
	"github.com/om00/golang-ecommerce/Token"
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

func (ctrl *Controller) Signup(w http.ResponseWriter, r *http.Request) {
	/*var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()*/

	var user Models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validation := validate.Struct(user)
	if validation != nil {
		http.Error(w, validation.Error(), http.StatusBadRequest)
		return
	}

	emailCnt, PhCnt, err := ctrl.db.UserAlreadyExist(user.Email, user.Phone)
	if err != nil {
		http.Error(w, "error while fetching data from db", http.StatusInternalServerError)
		return
	}
	if emailCnt > 0 && PhCnt > 0 {
		http.Error(w, "Both email or phone already registerd", http.StatusBadRequest)
		return
	} else if emailCnt > 0 && PhCnt == 0 {
		http.Error(w, "email alredy registerd", http.StatusBadRequest)
		return
	} else if emailCnt == 0 && PhCnt > 0 {
		http.Error(w, "phone number already registerd", http.StatusBadRequest)
		return
	}

	user.Password = HashPassword(user.Password)
	user.Token, user.Refresh_Token, err = Token.TokenGenrator(user.Email, user.First_Name, user.Last_Name, user.Phone)
	if err != nil {
		http.Error(w, "token generation fialed", http.StatusBadRequest)
		return
	}
	user.UserCart = make([]Models.ProductUser, 0)
	user.Address_Details = make([]Models.Address, 0)
	user.Order_Status = make([]Models.Order, 0)

	_, err = ctrl.db.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "you Singed up successfully!")

}

func (ctrl *Controller) Login(w http.ResponseWriter, r *http.Request) {

	var user Models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("error whe decode request, error = ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	userDb, err := ctrl.db.LoginUser(user.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error while fetching the data from db", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err == sql.ErrNoRows {
		log.Println("user does not exist")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	passwordisvalid, msg := verifyPassword(user.Password, userDb.Password)

	if !passwordisvalid {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, msg)
		return
	}

	user.Token, user.Refresh_Token, _ = Token.TokenGenrator(userDb.Email, userDb.First_Name, userDb.Last_Name, userDb.Phone)

	err = ctrl.db.UpdateToken(userDb.ID, user.Token, user.Refresh_Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("error while updatinng error=", err.Error())
	}

	w.WriteHeader(http.StatusFound)
	fmt.Fprint(w, "success")

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
// func (ctrl *Controller) SearchProductByQuery(w http.ResponseWriter, r *http.Request) {
// 	var searchProduct []Models.Product
// 	req_values := r.URL.Query()

// 	name_filter := req_values.Get("name")

// 	query := "SELECT id,productName,price,rating,image FROM Product WHERE productName=?"

// 	rows, err := db.Query(query, name_filter)

// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	for rows.Next() {
// 		var row_data Models.Product
// 		err := rows.Scan(&row_data.ID, &row_data.Product_Name, &row_data.Price, &row_data.Rating, &row_data.Image)
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		searchProduct = append(searchProduct, row_data)

// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	err = json.NewEncoder(w).Encode(searchProduct)

// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// }

type Controller struct {
	db *Database.DB
	//can dd more
}

func NewController(db *Database.DB) *Controller {
	return &Controller{
		db: db,
	}
}
