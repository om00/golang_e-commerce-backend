package Controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

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

func Signup(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

	var count_id, count_email, count_phone int64
	query := "SELECT COUNT(email) AS email,COUNT(phone) as phone From User WHERE email=? or phone=?"
	row, err := db.Query(query, user.Email, user.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Panic(err)

	}

	err = row.Scan(&count_id, &count_email, count_phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	if count_email > 0 && count_phone > 0 {
		http.Error(w, "Both email or phone already registerd", http.StatusBadRequest)
		return
	} else if count_email > 0 && count_phone == 0 {
		http.Error(w, "email alredy registerd", http.StatusBadRequest)
		return
	} else if count_email == 0 && count_phone > 0 {
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

	query = "INSERT INTO User COLOUNMS (`firstName`,`lastName`,`password`,`email`,`phone`,`token`,`refreshToken`,`userId`,`userCart`,`address`,`orderStatus`,`created_at`,`update_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"

	_, err = db.Exec(query, user.First_Name, user.Last_Name, user.Password, user.Email, user.Phone, user.Token, user.Refresh_Token, user.User_Id, user.UserCart, user.Address_Details, user.Order_Status, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "you Singed up successfully!")

}

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var user Models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	query := "SELECT id,firstName,lastName,phone,email,password FROM User where email=?"
	row := db.QueryRow(query, user.Email)
	/*if err!=nil{
		http.Error(w,err,http.StatusInternalServerError)
		log.println("error while executing query")
		return
	  }else*/if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "NO user exist with email %s")
		return
	}

	var firstname, lastname, phone, email, password string
	err = row.Scan(&user.ID, &firstname, &lastname, &phone, &email, &password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("error while executing query=%s", err)
		return
	}

	passwordisvalid, msg := verifyPassword(user.Password, password)

	if !passwordisvalid {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, msg)
		return
	}

	user.Token, user.Refresh_Token, _ = Token.TokenGenrator(email, firstname, lastname, phone)
	query = "UPDATE User SET token=?,refreshToken=? where id=?"

	_, err = db.Exec(query, user.Token, user.Refresh_Token, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error while updatinng error=%s", err)
	}

	w.WriteHeader(http.StatusFound)
	fmt.Fprint(w, "success")

}

func ProductViewAdmin(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

func SearchProduct(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var ProductList []Models.Product

	query := "Select id,productName,price,rating,image from Product"
	rows, err := db.Query(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var row Models.Product
		err := rows.Scan(&row.ID, &row.Product_Name, &row.Price, &row.Rating, &row.Image)
		if err != nil {
			fmt.Println("Error while scaning menthod")
			return
		}

		ProductList = append(ProductList, row)

	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(ProductList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func SearchProductByQuery(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var searchProduct []Models.Product
	req_values := r.URL.Query()

	name_filter := req_values.Get("name")

	query := "SELECT id,productName,price,rating,image FROM Product WHERE productName=?"

	rows, err := db.Query(query, name_filter)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var row_data Models.Product
		err := rows.Scan(&row_data.ID, &row_data.Product_Name, &row_data.Price, &row_data.Rating, &row_data.Image)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		searchProduct = append(searchProduct, row_data)

	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(searchProduct)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
