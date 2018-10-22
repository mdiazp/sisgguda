package app

import (
	"encoding/json"
	"net/http"
	"time"

	"gitlab.com/manuel.diaz/sisgguda/data"
)

type AddGroupHandler struct {
	app *App
}

func (h *AddGroupHandler) GetAccessRols() []string {
	return []string{RolAdmin}
}

func (h *AddGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var e error
	//Access Control
	author, ok := r.Context().Value("author").(data.User)
	if !ok {
		h.app.elogger.Println("Undefined author in context")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !containString(h.GetAccessRols(), author.Rol) {
		h.app.elogger.Println(http.StatusText(http.StatusForbidden))
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	//Reading JSON
	type Group struct{ Name, Description string }
	var newGroup Group
	e = ReadJsonFromRequest(r, &newGroup)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Validating not conflict with equal group names
	ok, e = h.app.datamodel.ExistGroupByName(newGroup.Name)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if ok {
		h.app.elogger.Println(e)
		http.Error(w, "Group name already exists.", http.StatusConflict)
		return
	}

	//Validating name
	if e := data.ValidGroupName(newGroup.Name); e != nil {
		h.app.elogger.Println(e)
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	//Validating description
	if e := data.ValidGroupDescription(newGroup.Description); e != nil {
		h.app.elogger.Println(e)
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	//Add Group to database
	group := data.Group{Name: newGroup.Name, Description: newGroup.Description, CreatedAt: time.Now()}
	e = h.app.datamodel.AddGroup(&group)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Output to JSON
	output, e := json.MarshalIndent(&group, "", "\t\t")
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Writing Response
	h.app.WriteResponse(w, http.StatusCreated, true, output)
	return
}
