package app

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/manuel.diaz/sisgguda/data"
)

type DeleteGroupHandler struct {
	app *App
}

func (h *DeleteGroupHandler) GetAccessRols() []string {
	return []string{RolAdmin}
}

func (h *DeleteGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	delGroupId, e := strconv.Atoi(vars["id"])
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	//Get user from database
	var group data.Group
	ok, e = h.app.datamodel.ExistGroupById(delGroupId)
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
	group, e = h.app.datamodel.GetGroupById(delGroupId)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Delete group in database
	e = h.app.datamodel.DeleteGroupById(group.Id)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Writing Response
	h.app.WriteResponse(w, http.StatusNoContent, true, nil)
	return
}
