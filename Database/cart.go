package Database

import (
	"fmt"
	"log"
	"context"

	"github.com/om00/golang-ecommerce/Models"
)

func AddProductToCart(ctx context.Context, app *Application,user_id,product_id int64) error{
	var user  Models.User
	var product  Models.Product
	query1:="SELECT id,userCart FROM User WHERE id=?"
	row,err:=app.db.QueryRow(query1,user_id)
	if err!=nil{
		log.Println("Error while fetching the data from table",err)
		return err
	}

	err:=row.Scan(&user.ID,&user.UserCart)

	if err!=nil{
		log.Println("Error while scanning the to struct",err)
		return err
	}

	query2:="SELECT * FROM Product WHERE id=?"
	row,err:=app.db.QueryRow(query2,product_id)

	if err !=nil{
		log.Println("Error while fetching the data from table",err)
		return err
	}

	err:=row.Scan(&product.ID,&product.Product_Name,&product.Price,&product.Rating,&product.Image)
	if err!=nil{
		log.Println("Error while scanning to the struct")
		return err
	}

	user.UserCart=append(user.UserCart,product)

	query3,err:=app.db.Prepare("UPDATE User SET userCart=? where id=?")
	if err!=nil{
		log.Println("Error while preparing update query",err)
		return err
	}
	defer query3.Close()

	row,err:=query3.Exec(user.UserCart,user.ID)

	if err!=nil{
		log.Println("Error while executing of update query",err)
		return err
	}
	defer row.Close()

	return nil
}




func RemoveItemFromCart(ctx context.Context,app *Application,user_id,product_id)error{
          var user Models.User
		  query1:="SELECT id,userCart FROM User WHERE id=?"
		  row,err:=app.db.QueryRow(query1,user_id)
		  if err!=nil {
			log.Println("err while executing the query",error)
			return err
		  }

		  err:=row.Scan(&user.ID,&user.UserCart)
		  if err!=nil{
			log.Println("Error while scaning value into the struct",error)
			return err
		  }

		  for i:=0;i<len(user.UserCart);i++{
			if product_id==user.userCart[i].id{
				user.UserCart=append(user.UserCart[:i],user.UserCart[i+1:])
				i--
			}

			}

		  query2,err:=app.db.Prepare("UPDATE User SET userCart=? WHERE id=?")

		  if err!=nil{
			log.Println("Error while preparing query for updation")
			return err
		  }

		  row,err:=query2.Exec(user.UserCart,user.ID)
          if err!=nil{
			log.Println("Error while executing query for updation",err)
			return err
		  }

		  return nil

		  }


func BuyIteamFromCart(ctx context.Context,app *Application,user_id int64) error{

	var order Models.Order
	var user  Models.User

	order.Order_Cart=make([]Models.ProductUser, 0)

	query:="SELECT id,userCart FROM User WHERE id=?"
	row,err:=app.db.QueryRow(query,user_id)

	if err!=nil{
		log.Println("Erorr while executing the error")
		return err
	}

	err:=row.Scan(&user.id,&user.UserCart)

	if err!=nil{
		log.Println("Error while scanning to struct")
		return err
	}
    
	for _,value :=range user.UserCart{
              order.Price=order.Price+value.Price
			  order.Order_Cart=append(order.Order_Cart,value)
	}

	order.Payment_Method.COD=true

	order.Discount=0
	
	query2,err:=app.db.Prepare("UPDATE User SET userCart=? WHERE id=?")

	if err!=nil{
		log.Println("Error while Preparing query")
		return err
	}

	res,err:=query2.Exec(nil,id)
    if err!=nil{
		log.Println("Error while updating the value",err)
		return err
	}

	query3,err:=app.db.Prepare("INSERT INTO Order COLOUMNS(orderlist,created_at,updated_at,price,discount,payment)
	VALUES  (?,?,?,?,?,?)")
	

    res,err:=query3.Exec(order.Order_Cart,tiPayment_Methodme.Now().Format("2006-01-02 15:04:05"),time.Now().Format("2006-01-02 15:04:05"),order.Price
     order.Discount,order.Payment_Method)

	 if err!=nil{
		log.Println("Error while inserting the data into the table",err)
		return err
	 }

	 return nil


                               
}

func  InstantBuyer(ctx context.Context,app *Application,product_id int64)error{
      
	  var product_details Models.Product

	  var order Models.Order

	  query="SELECT * FROM Product WHERE id=? "

	  row,err:=app.db.QueryRow(query,product_id)
      if err!=nil{
		log.Println("Error while exectuint the query",err)
		return err
	  }

	  err:=row.Scan(&product_details.ID,&product_details.Product_Name,&product_details.Price,
		&product_details.Image,&product_details.Rating)

	if err!=nil{
		log.Println("Error while scannnig the struct")
		return err
	}

	order.Price=product_details.Price
	order.Order_Cart=append(order.Order_Cart,product_details)

	order.Discount=0
	order.Payment_Method.COD=true


	query3,err:=app.db.Prepare("INSERT INTO Order COLOUMNS(orderlist,created_at,updated_at,price,discount,payment)
	VALUES  (?,?,?,?,?,?)")
	

    res,err:=query3.Exec(order.Order_Cart,tiPayment_Methodme.Now().Format("2006-01-02 15:04:05"),time.Now().Format("2006-01-02 15:04:05"),order.Price
     order.Discount,order.Payment_Method)

	 if err!=nil{
		log.Println("Error while inserting the data into the table",err)
		return err
	 }

	 return nil


}


