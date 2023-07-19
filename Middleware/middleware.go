package Middleware

import (
	"context"
	"log"
	"net/http"

	token "github.com/om00/golang-ecommerce/Token"
)

func Authentication(w http.ResponseWriter, r *http.Request) {
	Client_token := r.Header.Get("Token")

	if Client_token == "" {
		log.Println("No authtiection provide in header")
		http.Error(w, "no authietication provided in header", http.StatusInternalServerError)
		return
	}

	claims, err := token.ValidatedToekn(Client_token)

	if err != "" {
		log.Println("token is not correct")
		http.Error(w, "token is not correcty", http.StatusBadRequest)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), "email", claims.Email))

	// Call the next handler in the chain
	next(w, r)

}
