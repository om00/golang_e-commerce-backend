package Database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")

	go func() {
		db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName))
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

func (db *DB) RunMigrations() error {
	driver, err := mysql.WithInstance(db.mainDB.DB, &mysql.Config{})
	if err != nil {
		log.Println("Error to get Db instance", err.Error())
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://Database/Migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Println("Error to get migration instance", err.Error())
		return err
	}

	if len(os.Args) > 3 {
		log.Println("More argument is provided  then expected in commandline")
		return fmt.Errorf("more argument is provided  then expected in commandline")
	}

	// Read the command-line arguments
	if len(os.Args) > 2 {

		command := os.Args[1]
		action := ""
		if len(os.Args) > 2 {
			action = os.Args[2]
		}

		// If the first argument is "withMigration", decide whether to run "up" or "down"
		if command == "withMigration" {
			var msg string
			switch strings.ToLower(action) {
			case "up":
				msg = "run the Migration"
				err = m.Up()
			case "down":
				msg = "rollabck the Migration"
				err = m.Steps(-1)

			}

			if err != nil && err != migrate.ErrNoChange {
				log.Println("error while"+msg, err)
				return fmt.Errorf("error while %s, error= %s", msg, err.Error())
			}
		}
	}

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println("Error while running migration", err.Error())
		return err
	}

	return nil
}
