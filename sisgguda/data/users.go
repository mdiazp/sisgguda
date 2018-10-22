package data

import (
	"errors"
	"time"
)

type User struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Rol         string `json:"rol"`
	CreatedAt   time.Time
}

const (
	UserDescriptionMaxSize = 300
)

func (model *DataModel) ExistUserById(id int) (ok bool, e error) {
	e = model.Db.QueryRow("select exists (select * from users where id = $1)", id).Scan(&ok)
	return
}

func (model *DataModel) ExistUserByUsername(username string) (ok bool, e error) {
	e = model.Db.QueryRow("select exists (select * from users where username = $1)", username).Scan(&ok)
	return
}
func (model *DataModel) GetUserById(id int) (user User, err error) {
	user = User{}
	err = model.Db.QueryRow("select * from users where id = $1", id).Scan(&user.Id, &user.Username, &user.Description, &user.Rol, &user.CreatedAt)
	return
}

func (model *DataModel) GetUserByUsername(username string) (user User, err error) {
	user = User{}
	err = model.Db.QueryRow("select * from users where username = $1", username).Scan(&user.Id, &user.Username, &user.Description, &user.Rol, &user.CreatedAt)
	return
}

func (model *DataModel) AddUser(user *User) (err error) {
	statement := "insert into users (username, description, rol, createdAt ) values ($1, $2, $3, $4) returning id"
	stmt, err := model.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(user.Username, user.Description, user.Rol, user.CreatedAt).Scan(&user.Id)

	return
}

func (model *DataModel) UpdateUser(user *User) (err error) {
	statement := "update users SET username=$1, description=$2, rol=$3, createdAt=$4 where id = $5"
	stmt, err := model.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Username, user.Description, user.Rol, user.CreatedAt, user.Id)
	return
}

func (model *DataModel) DeleteUserById(id int) (err error) {
	statement := "delete from users where id = $1"
	stmt, err := model.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return
}

func ValidUserDescription(description string) (e error) {
	if len(description) > UserDescriptionMaxSize {
		e = errors.New("User description is too long.")
	}
	return
}
