package app

import (
	"net/http"

	"gitlab.com/manuel.diaz/sisgguda/data"
)

type DeleteGroupSpecialistHandler struct {
	app *App
}

func (h *DeleteGroupSpecialistHandler) GetAccessRols() []string {
	return []string{RolAdmin}
}

func (h *DeleteGroupSpecialistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	type GroupSpecialist struct{ GroupId, UserId int }
	var newSpecialist GroupSpecialist
	e = ReadJsonFromRequest(r, &newSpecialist)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Validating group existence
	ok, e = h.app.datamodel.ExistGroupById(newSpecialist.GroupId)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !ok {
		h.app.elogger.Println(e)
		http.Error(w, "Group Not Found", http.StatusBadRequest)
		return
	}

	//Get group from database
	group, e := h.app.datamodel.GetGroupById(newSpecialist.GroupId)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Validating user existence
	ok, e = h.app.datamodel.ExistUserById(newSpecialist.UserId)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !ok {
		h.app.elogger.Println(e)
		http.Error(w, "User Not Found", http.StatusBadRequest)
		return
	}

	//Get user from database
	user, e := h.app.datamodel.GetUserById(newSpecialist.UserId)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Validating specialist exists
	ok, e = h.app.datamodel.ExistGroupSpecilist(&group, &user)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !ok {
		h.app.elogger.Println(e)
		http.Error(w, "This group dont have this user as specialist.", http.StatusConflict)
		return
	}

	//Delete Specialist in database
	e = h.app.datamodel.DeleteGroupSpecialist(&group, &user)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Writing Response
	h.app.WriteResponse(w, http.StatusNoContent, true, nil)
	return
}
