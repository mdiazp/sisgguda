package data

type GroupAdUser struct {
	GroupId    string `json:"groupId"`
	AdUsername string `json:"adUsername"`
}

func (model *DataModel) ExistGroupAdUser(group *Group, username string) (ok bool, e error) {
	e = model.Db.QueryRow("select exists (select * from group_adusers where group_id = $1 and ad_username = $2)", group.Id, username).Scan(&ok)
	return
}

func (model *DataModel) AddGroupAdUser(group *Group, username string) (err error) {
	statement := "insert into group_adusers (group_id, ad_username ) values ($1, $2)"
	stmt, err := model.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(group.Id, username)
	return
}

func (model *DataModel) DeleteGroupAdUser(group *Group, username string) (err error) {
	statement := "delete from group_adusers where group_id = $1 and ad_username = $2"
	stmt, err := model.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(group.Id, username)
	return
}
