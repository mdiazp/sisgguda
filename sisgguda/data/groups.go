package data

import (
	"errors"
	"time"
)

type Group struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   time.Time
}

const (
	GroupNameMaxSize        = 100
	GroupDescriptionMaxSize = 300
)

func (model *DataModel) ExistGroupById(id int) (ok bool, e error) {
	e = model.Db.QueryRow("select exists (select * from groups where id = $1)", id).Scan(&ok)
	return
}

func (model *DataModel) ExistGroupByName(name string) (ok bool, e error) {
	e = model.Db.QueryRow("select exists (select * from groups where name = $1)", name).Scan(&ok)
	return
}
func (model *DataModel) GetGroupById(id int) (group Group, err error) {
	group = Group{}
	err = model.Db.QueryRow("select * from groups where id = $1", id).Scan(&group.Id, &group.Name, &group.Description, &group.CreatedAt)
	return
}

func (model *DataModel) GetGroupByName(name string) (group Group, err error) {
	group = Group{}
	err = model.Db.QueryRow("select * from groups where name = $1", name).Scan(&group.Id, &group.Name, &group.Description, &group.CreatedAt)
	return
}

func (model *DataModel) AddGroup(group *Group) (err error) {
	statement := "insert into groups (name, description, createdAt ) values ($1, $2, $3) returning id"
	stmt, err := model.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(group.Name, group.Description, group.CreatedAt).Scan(&group.Id)

	return
}

func (model *DataModel) UpdateGroup(group *Group) (err error) {
	statement := "update groups SET name=$1, description=$2, createdAt=$3 where id = $4"
	stmt, err := model.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(group.Name, group.Description, group.CreatedAt, group.Id)
	return
}

func (model *DataModel) DeleteGroupById(id int) (err error) {
	statement := "delete from groups where id = $1"
	stmt, err := model.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return
}

func ValidGroupName(name string) (e error) {
	nameLen := len(name)
	if nameLen == 0 {
		e = errors.New("Group name can not be empty")
	}
	if nameLen > GroupNameMaxSize {
		e = errors.New("Group name is too long.")
	}
	return
}

func ValidGroupDescription(description string) (e error) {
	if len(description) > GroupDescriptionMaxSize {
		e = errors.New("Group description is too long.")
	}
	return
}
