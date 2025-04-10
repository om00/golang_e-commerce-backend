# golang_e-commerce-backend

It is simple e-commerce beckend build in golang.
it follows mvc structue.

Routes Folder - it contains all entry points for api. 
Controller Folder- All logic is written in controller.
Database Folder- all logic to interact with database is written here.
                 it has migration folder as well so new tables can be added.


setup -
1.golang should be installed on local system.
2. mysql server should  be installed  on local system .

modify env according to your database ,user, and ports 
you can choose on which port service will just need to in env

to run start server command  -
                      go run main.go

if want to start with migration 
                    go run main.go withMigration up

if need to rollback the migration 
                    go run main.go withMigration down



next improvments -
// make it continer based using docker 
