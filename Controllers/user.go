package Controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/om00/golang-ecommerce/Models"
	"github.com/om00/golang-ecommerce/Token"
)

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
		log.Printf("error where getting the data from db %s", err)
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
	user.Token, user.Refresh_Token, err = Token.TokenGenrator(user.ID, user.Email, user.First_Name, user.Last_Name, user.Phone)
	if err != nil {
		http.Error(w, "token generation fialed", http.StatusBadRequest)
		return
	}
	user.UserCart = make([]Models.ProductUser, 0)
	user.Address_Details = make([]Models.Address, 0)
	user.Order_Status = make([]Models.Order, 0)

	_, err = ctrl.db.CreateUser(user)
	if err != nil {
		log.Printf("error where creating the user error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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

	userDb, err := ctrl.db.GetUser(0, user.Email)
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

	user.Token, user.Refresh_Token, _ = Token.TokenGenrator(userDb.ID, userDb.Email, userDb.First_Name, userDb.Last_Name, userDb.Phone)

	err = ctrl.db.UpdateToken(userDb.ID, user.Token, user.Refresh_Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("error while updatinng error=", err.Error())
	}

	response := map[string]interface{}{
		"success":       true,
		"token":         user.Token,
		"refresh_token": user.Refresh_Token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (ctrl *Controller) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req Models.RefreshTokenReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Token == "" || req.UserId == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	claime, errmsg := Token.ValidatedToekn(req.Token)
	if errmsg != "" {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	user, err := ctrl.db.GetUser(claime.Id, "")
	if err != nil {
		http.Error(w, "Error feching user infromations", http.StatusInternalServerError)
		return
	}

	// Generate new access token
	newAccessToken, newRefreshToken, err := Token.TokenGenrator(user.ID, user.Email, user.First_Name, user.Last_Name, user.Phone)
	if err != nil {
		http.Error(w, "Error generating new tokens", http.StatusInternalServerError)
		return
	}

	err = ctrl.db.UpdateToken(user.ID, newAccessToken, newRefreshToken)
	if err != nil {
		http.Error(w, "Error in updating  tokens", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"accessToken":  newAccessToken,
		"refreshToken": newRefreshToken,
	})
}
