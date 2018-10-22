package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/manuel.diaz/sisgguda/data"
)

type UpdateGroupHandler struct {
	app *App
}

func (h *UpdateGroupHandler) GetAccessRols() []string {
	return []string{RolAdmin}
}

func (h *UpdateGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	//Get groupId from url
	vars := mux.Vars(r)
	updGroupId, e := strconv.Atoi(vars["id"])
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	//Get group from database
	var group data.Group
	ok, e = h.app.datamodel.ExistGroupById(updGroupId)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !ok {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	group, e = h.app.datamodel.GetGroupById(updGroupId)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Reading JSON
	type Group struct{ Name, Description string }
	var update Group
	e = ReadJsonFromRequest(r, &update)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Validating Name
	if e := data.ValidGroupName(update.Name); e != nil {
		h.app.elogger.Println(e)
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	//Validating description
	if e := data.ValidGroupDescription(update.Description); e != nil {
		h.app.elogger.Println(e)
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	//Update Group in database
	group.Name = update.Name
	group.Description = update.Description
	e = h.app.datamodel.UpdateGroup(&group)

	//Output to JSON
	var output []byte
	if e == nil {
		output, e = json.MarshalIndent(&group, "", "\t\t")
	}

	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Writing Response
	h.app.WriteResponse(w, 200, true, output)
	return
}
