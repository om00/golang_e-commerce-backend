package Controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/om00/golang-ecommerce/Models"
)

func AddAddress(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	var address Models.Address
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		log.Println("err while decoidng request json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func DeleteAddress(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	req_value := r.URL.Query()
	user_id := req_value.Get("id")

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	address := make([]Models.Address, 0)

	query := "UPDATE User SET address=? WHERE id=?"
	_, err := DB.Exec(query, address, user_id)

	if err != nil {
		log.Println("error while deleting the address")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted address")

}

func EditAddress(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

}
