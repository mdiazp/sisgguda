package app

import (
	"net/http"

	"gitlab.com/manuel.diaz/sisgguda/data"
)

type AddGroupSpecialistHandler struct {
	app *App
}

func (h *AddGroupSpecialistHandler) GetAccessRols() []string {
	return []string{RolAdmin}
}

func (h *AddGroupSpecialistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	type GroupSpecialist struct {
		GroupId int
		UserId  int
	}
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
		h.app.elogger.Println("Group Not Found, groupId =", newSpecialist.GroupId)
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
		h.app.elogger.Println("User Not Found, userId =", newSpecialist.UserId)
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

	//Validating specialist not exists
	ok, e = h.app.datamodel.ExistGroupSpecilist(&group, &user)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if ok {
		h.app.elogger.Println(e)
		http.Error(w, "This group already have this user as specialist.", http.StatusConflict)
		return
	}

	//Add Specialist to database
	e = h.app.datamodel.AddGroupSpecialist(&group, &user)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Writing Response
	h.app.WriteResponse(w, http.StatusCreated, true, nil)
	return
}
