package Database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	timeout := 10 * time.Second

	ch := make(chan *sql.DB, 1)
	errch := make(chan error, 1)

	go func() {
		db, err := sql.Open("mysql", "root:go-lang-ec@tcp(127.0.0.1:3307)/")
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
