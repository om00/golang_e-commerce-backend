package Database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/om00/golang-ecommerce/Models"
)

func (db *DB) AddProductToCart(ctx context.Context, user_id, product_id int64) error {
	var user Models.User
	var product Models.ProductUser
	query1 := "SELECT id,userCart FROM User WHERE id=?"
	row := db.mainDB.QueryRow(query1, user_id)

	err := row.Scan(&user.ID, &user.UserCart)

	if err != nil {
		log.Println("Error while scanning the to struct", err)
		return err
	}

	if err == sql.ErrNoRows {
		log.Println("No rows found.")
		return err
	}

	query2 := "SELECT * FROM Product WHERE id=?"
	row = db.mainDB.QueryRow(query2, product_id)

	err = row.Scan(&product.ID, &product.Product_Name, &product.Price, &product.Rating, &product.Image)
	if err != nil {
		log.Println("Error while scanning to the struct")
		return err
	}

	if err == sql.ErrNoRows {
		log.Println("No rows found.")
		return err
	}

	user.UserCart = append(user.UserCart, product)

	query3, err := db.mainDB.Prepare("UPDATE User SET userCart=? where id=?")
	if err != nil {
		log.Println("Error while preparing update query", err)
		return err
	}
	defer query3.Close()

	_, err = query3.Exec(user.UserCart, user.ID)

	if err != nil {
		log.Println("Error while executing of update query", err)
		return err
	}

	return nil
}

func (db *DB) RemoveItemFromCart(ctx context.Context, user_id, product_id int64) error {
	var user Models.User
	query1 := "SELECT id,userCart FROM User WHERE id=?"
	row := db.mainDB.QueryRow(query1, user_id)

	err := row.Scan(&user.ID, &user.UserCart)
	if err != nil {
		log.Println("Error while scaning value into the struct", err)
		return err
	}

	if err == sql.ErrNoRows {
		log.Println("No rows found.")
		return err
	}

	for i := 0; i < len(user.UserCart); i++ {
		if product_id == user.UserCart[i].ID {
			user.UserCart = append(user.UserCart[:i], user.UserCart[i+1:]...)
			i--
		}

	}

	query2, err := db.mainDB.Prepare("UPDATE User SET userCart=? WHERE id=?")

	if err != nil {
		log.Println("Error while preparing query for updation")
		return err
	}

	_, err = query2.Exec(user.UserCart, user.ID)
	if err != nil {
		log.Println("Error while executing query for updation", err)
		return err
	}

	return nil

}

func (db *DB) BuyIteamFromCart(ctx context.Context, user_id int64) error {

	var order Models.Order
	var user Models.User

	order.Order_Cart = make([]Models.ProductUser, 0)

	query := "SELECT id,userCart FROM User WHERE id=?"
	row := db.mainDB.QueryRow(query, user_id)

	err := row.Scan(&user.ID, &user.UserCart)

	if err == sql.ErrNoRows {
		log.Println("No rows found.")
		return err
	}

	if err != nil {
		log.Println("Error while scanning to struct")
		return err
	}

	for _, value := range user.UserCart {
		order.Price = order.Price + value.Price
		order.Order_Cart = append(order.Order_Cart, value)
	}

	order.Payment_Method.COD = true

	order.Discount.Valid = false

	query2, err := db.mainDB.Prepare("UPDATE User SET userCart=? WHERE id=?")

	if err != nil {
		log.Println("Error while Preparing query")
		return err
	}

	_, err = query2.Exec(nil, user_id)
	if err != nil {
		log.Println("Error while updating the value", err)
		return err
	}

	query3, err := db.mainDB.Prepare("INSERT INTO Order COLOUMNS(orderlist,created_at,updated_at,price,discount,payment) VALUES  (?,?,?,?,?,?)")

	_, err = query3.Exec(order.Order_Cart, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), order.Price, order.Discount, order.Payment_Method.COD)

	if err != nil {
		log.Println("Error while inserting the data into the table", err)
		return err
	}

	return nil

}

func (db *DB) InstantBuyer(ctx context.Context, product_id int64) error {

	var product_details Models.ProductUser

	var order Models.Order

	query := "SELECT * FROM Product WHERE id=? "

	row := db.mainDB.QueryRow(query, product_id)

	err := row.Scan(&product_details.ID, &product_details.Product_Name, &product_details.Price,
		&product_details.Image, &product_details.Rating)

	if err != nil {
		log.Println("Error while scannnig the struct")
		return err
	}

	if err == sql.ErrNoRows {
		log.Println("No rows found.")
		return err
	}

	order.Price = product_details.Price
	order.Order_Cart = append(order.Order_Cart, product_details)

	order.Discount.Valid = false
	order.Payment_Method.COD = true

	query3, err := db.mainDB.Prepare("INSERT INTO Order COLOUMNS(orderlist,created_at,updated_at,price,discount,payment) VALUES  (?,?,?,?,?,?)")

	_, err = query3.Exec(order.Order_Cart, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), order.Price,
		order.Discount, order.Payment_Method.COD)

	if err != nil {
		log.Println("Error while inserting the data into the table", err)
		return err
	}

	return nil

}
