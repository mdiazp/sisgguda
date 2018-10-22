package app

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/manuel.diaz/sisgguda/data"
)

type DeleteUserHandler struct {
	app *App
}

func (h *DeleteUserHandler) GetAccessRols() []string {
	return []string{RolSuperAdmin, RolAdmin}
}

func (h *DeleteUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	//Get userId from url
	vars := mux.Vars(r)
	delUserId, e := strconv.Atoi(vars["id"])
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	//Get user from database
	var user data.User
	ok, e = h.app.datamodel.ExistUserById(delUserId)
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
	user, e = h.app.datamodel.GetUserById(delUserId)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//SuperAdmin never can be deleted
	if user.Rol == RolSuperAdmin {
		h.app.elogger.Println("SuperAdmin users never can be deleted")
		http.Error(w, "SuperAdmin users never can be deleted.", http.StatusForbidden)
		return
	}

	//Only SuperAdmin can delete admin users
	if user.Rol == RolAdmin && author.Rol != RolSuperAdmin {
		h.app.elogger.Println("Only SuperAdmin can delete admin users")
		http.Error(w, "Only SuperAdmin can delete admin users.", http.StatusForbidden)
		return
	}

	//Delete user in database
	e = h.app.datamodel.DeleteUserById(user.Id)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Writing Response
	h.app.WriteResponse(w, 204, true, nil)
	return
}
