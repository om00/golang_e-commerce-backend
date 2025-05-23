package Models

import (
	"database/sql"
	"time"
)

type User struct {
	ID              int64         `db:"id" json:"id" `
	First_Name      string        `db:"firstName" json:"firstName" validate:"required"`
	Last_Name       string        `db:"lastName" json:"lastName" validate:"required"`
	Password        string        `db:"password" json:"password" validate:"required"`
	Email           string        `db:"email" json:"email" validate:"required"`
	Phone           string        `db:"phone" json:"phone" validate:"required"`
	Token           string        `db:"token"`
	Refresh_Token   string        `db:"refershToken"`
	Created_At      time.Time     `db:"created_at"`
	Upadated_At     time.Time     `db:"updated_at"`
	User_Id         string        `db:"userId"`
	UserCart        []ProductUser `db:"userCart"`
	Address_Details []Address     `db:"Address"`
	Order_Status    []Order       `db:"orderStatus"`
}

type Product struct {
	ID           int64           `db:"id"`
	Product_Name string          `db:"productName"`
	Price        float64         `db:"price"`
	Rating       sql.NullFloat64 `db:"rating"`
	Image        string          `db:"image"`
}

type ProductUser struct {
	ID           int64
	Product_Name string
	Price        float64
	Rating       sql.NullFloat64
	Image        string
}

type Address struct {
	ID      int64
	House   string
	Street  string
	City    string
	Pincode string
}
type Order struct {
	ID             int64           `db:"id"`
	Order_Cart     []ProductUser   `db:"orderList"`
	Created_At     time.Time       `db:"created_at"`
	Updated_at     time.Time       `db:"updated_at"`
	Price          float64         `db:"price"`
	Discount       sql.NullFloat64 `db:"discount"`
	Payment_Method Payment         `db:"payment"`
}

type Payment struct {
	Digital bool
	COD     bool
}

type ProductQuery struct {
	Name  string
	Price float64
}

type RefreshTokenReq struct {
	Token  string `json:"token"`
	UserId int64  `josn:"userid"`
}
