package Database

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	mainDB *sqlx.DB
	//can add if need another attributes
}

func InitDB() (*sqlx.DB, error) {
	timeout := 10 * time.Second

	ch := make(chan *sqlx.DB, 1)
	errch := make(chan error, 1)

	go func() {
		db, err := sqlx.Open("mysql", "Om:Om07@Golang@tcp(127.0.0.1:3306)/GoEmcDB")
		if err != nil {
			errch <- err
			return
		}

		err = db.Ping()
		if err != nil {
			errch <- err
			return
		}

		ch <- db
	}()

	select {
	case db := <-ch:
		fmt.Println("connection is established sucessfully to the database")
		return db, nil
	case err := <-errch:
		fmt.Println("Error while connecting to the database")
		return nil, err
	case <-time.After(timeout):
		fmt.Println("timeout")
		return nil, errors.New("taking too much time to connect")

	}

}

func NewDB() (*DB, error) {
	db, err := InitDB()
	if err != nil {
		fmt.Println("error while Establish connection")
		return nil, err
	}
	return &DB{mainDB: db}, nil
}
