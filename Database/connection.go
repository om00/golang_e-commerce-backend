package Database

import(
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
)




func initDB() *sql.DB,error{
    timeout := 10 * time.Second

	ch:=make(chan *sql.DB,1)
	errch:=make(chan error,1 )

	go func(){
	db,err:=sql.open("mysql","username:password@tcp(host:port)/database");
	if err!=nil{
		errch<-err
		return 
	}

	err:=db.Ping()
	if err!=nil{
		errch<-err
		return
	}

    ch<-db
}()


   select {
   case db:=<-ch
	    fmt.Println("connection is established sucessfully to the database")
		return db,nil
   case err:=<-errch:
	    log.Println("Error while connecting to the database")
		return nil,err
    case <-After.timeout(timeout)
	    log.Println("timeout")
	    return nil,log.Errorf("taking too much time to connect")

   }

}


var Db *sql.DB =initDB()
