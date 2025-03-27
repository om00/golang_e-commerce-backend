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

	err = row.Scan(&emailCount, &phCount)
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

func (db *DB) LoginUser(email string) (Models.User, error) {

	query, args, err := squirrel.Select("id", "firstName", "lastName", "email", "password").
		From("User").
		Where(squirrel.Eq{"email": email}).
		Limit(1).ToSql()

	if err != nil {
		return Models.User{}, err
	}

	var user Models.User
	err = db.mainDB.Get(&user, query, args...)
	if err != nil {
		return Models.User{}, err
	}

	return user, nil
}

func (db *DB) UpdateToken(id int64, token, refreshToken string) (err error) {
	query, args, err := squirrel.Update("User").
		Set("token", token).
		Set("refreshToken", refreshToken).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return err
	}
	_, err = db.mainDB.Exec(query, args)
	if err != nil {
		return nil
	}

	return nil
}
