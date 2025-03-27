package Database

import (
	"github.com/Masterminds/squirrel"
	"github.com/om00/golang-ecommerce/Models"
)

func (db *DB) UserAlreadyExist(email string, phone string) (emailCount, phCount int16, err error) {

	query := "SELECT COUNT(email) AS email,COUNT(phone) as phone From User WHERE email=? or phone=?"
	row, err := db.mainDB.Query(query, email, phone)
	if err != nil {
		return 0, 0, nil

	}

	err = row.Scan(&emailCount, phCount)
	if err != nil {
		return 0, 0, err
	}

	return emailCount, phCount, nil
}

func (db *DB) CreateUser(user Models.User) (ID int64, err error) {

	queryBuilder := squirrel.Insert("User").
		Columns(
			"firstName",
			"lastName",
			"password",
			"email",
			"phone",
			"token").
		Values(
			user.First_Name,
			user.Last_Name,
			user.Password,
			user.Email,
			user.Phone,
			user.Token,
			user.Refresh_Token,
		)
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, err
	}
	result, err := db.mainDB.Exec(query, args)
	if err != nil {
		return 0, nil

	}

	insertID, _ := result.LastInsertId()

	return insertID, nil
}
