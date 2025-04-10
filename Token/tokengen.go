package Token

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type SingedUpDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Phone      string
	Id         int64
	jwt.StandardClaims
}

var SECRET_KEY string

func TokenGenrator(userId int64, email string, first_name string, last_name string, phone string) (SignedToken, SignedRefreshToekn string, err error) {
	claime := &SingedUpDetails{
		Email:      email,
		First_Name: first_name,
		Last_Name:  last_name,
		Phone:      phone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		}}

	Refresh_claime := &SingedUpDetails{
		Id: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		}}

	SignedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claime).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println("Error while creating tokens", err)
		return "", "", err
	}

	SignedRefreshToekn, err = jwt.NewWithClaims(jwt.SigningMethodHS256, Refresh_claime).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Println("Error while creating tokens", err)
		return "", "", err
	}

	return SignedToken, SignedRefreshToekn, nil
}

func ValidatedToekn(singedtoken string) (claims *SingedUpDetails, errMsg string) {
	token, err := jwt.ParseWithClaims(singedtoken, &SingedUpDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		errMsg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SingedUpDetails)

	if !ok {
		errMsg = "the token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		errMsg = "token is already expired"
		return
	}

	return claims, errMsg
}
