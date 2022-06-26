package models

import (
	"database/sql"
	"davidwah/login/config"
	"davidwah/login/entities"
)

type UserModel struct {
	db *sql.DB
}

func NewUserModel() *UserModel {
	conn, err := config.DBConn()

	if err != nil {
		panic(err)
	}

	return &UserModel{
		db: conn,
	}
}

func (u UserModel) Where(user *entities.User, fileldName, fieldValue string) error {

	row, err := u.db.Query("select id, nama, email, username, password from users where "+fileldName+" = ? limit 1", fieldValue)

	if err != nil {
		return err
	}
	defer row.Close()

	for row.Next() {
		row.Scan(&user.Id, &user.Nama, &user.Email, &user.Username, &user.Password)
	}

	return nil

}
