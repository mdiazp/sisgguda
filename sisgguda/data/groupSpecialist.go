package data

type GroupSpecialist struct {
	GroupId string `json:"groupId"`
	UserId  string `json:"userId"`
}

func (model *DataModel) ExistGroupSpecilist(group *Group, user *User) (ok bool, e error) {
	e = model.Db.QueryRow("select exists (select * from group_specialists where group_id = $1 and user_id = $2)", group.Id, user.Id).Scan(&ok)
	return
}

func (model *DataModel) AddGroupSpecialist(group *Group, user *User) (err error) {
	statement := "insert into group_specialists (group_id, user_id ) values ($1, $2)"
	stmt, err := model.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(group.Id, user.Id)
	return
}

func (model *DataModel) DeleteGroupSpecialist(group *Group, user *User) (err error) {
	statement := "delete from group_specialists where group_id = $1 and user_id = $2"
	stmt, err := model.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(group.Id, user.Id)
	return
}
