package Token

import (
	"database/sql"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type SingedUpDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Phone      string
	jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenrator(email string, first_name string, last_name string, phone string) (SignedToken, SignedRefreshToekn string, err error) {
	claime := &SingedUpDetails{
		Email:      email,
		First_Name: first_name,
		Last_Name:  last_name,
		Phone:      phone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		}}

	Refresh_claime := &SingedUpDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		}}

	SignedToken, err = jwt.NewWithClaims(jwt.SigningMethodES256, claime).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Println("Error while creating tokens", err)
		return "", "", err
	}

	SignedRefreshToekn, err = jwt.NewWithClaims(jwt.SigningMethodES256, Refresh_claime).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Println("Error while creating tokens", err)
		return "", "", err
	}

	return SignedToken, SignedRefreshToekn, nil
}

func ValidatedToekn(singedtoken string) (claims *SingedUpDetails, msg string) {
	token, err := jwt.ParseWithClaims(singedtoken, &SingedUpDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SingedUpDetails)

	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is already expired"
		return
	}

	return claims, msg

}

func UpadateAllTokens(signedtoken, refresh_singedtoken string, user_id int64, db *sql.DB) {

	/*ctx, canecl := context.WithTimeout(context.Background(), 100*time.Second)*/

	query, err := db.Prepare("UPDATE User SET token=?,refreshToken=?,updated_at=? where id=?")
	if err != nil {
		log.Println("error while preparing the query")
		return
	}

	_, err = query.Exec(signedtoken, refresh_singedtoken, time.Now().Format("2006-01-02 15:04:05"), user_id)
	/*defer cancel()*/

	if err != nil {
		log.Println("error while executing  the  query")
		return
	}

}
