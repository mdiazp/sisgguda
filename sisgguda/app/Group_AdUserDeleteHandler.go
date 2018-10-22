package app

import (
	"net/http"

	"gitlab.com/manuel.diaz/sisgguda/data"
)

type DeleteGroupAdUserHandler struct {
	app *App
}

func (h *DeleteGroupAdUserHandler) GetAccessRols() []string {
	return []string{RolAdmin}
}

func (h *DeleteGroupAdUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	type GroupAdUser struct {
		GroupId    int
		AdUsername string
	}
	var newAdUser GroupAdUser
	e = ReadJsonFromRequest(r, &newAdUser)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Validating group existence
	ok, e = h.app.datamodel.ExistGroupById(newAdUser.GroupId)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !ok {
		h.app.elogger.Println("Group Not Found, groupId =", newAdUser.GroupId)
		http.Error(w, "Group Not Found", http.StatusBadRequest)
		return
	}

	//Get group from database
	group, e := h.app.datamodel.GetGroupById(newAdUser.GroupId)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Validating user in ldap server
	if e := h.app.ValidAdUser(newAdUser.AdUsername); e != nil {
		h.app.elogger.Println(e)
		http.Error(w, "Invalid user in ldap server.", http.StatusBadRequest)
		return
	}

	//Validating aduser exists in group
	ok, e = h.app.datamodel.ExistGroupAdUser(&group, newAdUser.AdUsername)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !ok {
		h.app.elogger.Println(e)
		http.Error(w, "This group dont have this user", http.StatusConflict)
		return
	}

	//Validating author be specialist of group, or author be admin
	if author.Rol != RolAdmin {
		ok, e = h.app.datamodel.ExistGroupSpecilist(&group, &author)
		if e != nil {
			h.app.elogger.Println(e)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !ok {
			h.app.elogger.Println(e)
			http.Error(w, "You are not specialist of this group.s", http.StatusForbidden)
			return
		}
	}

	//Delete GroupAdUser to database
	e = h.app.datamodel.DeleteGroupAdUser(&group, newAdUser.AdUsername)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Writing Response
	h.app.WriteResponse(w, http.StatusNoContent, true, nil)
	return
}
