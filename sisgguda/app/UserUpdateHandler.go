package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"gitlab.com/manuel.diaz/sisgguda/data"
)

type UpdateUserHandler struct {
	app *App
}

func (h *UpdateUserHandler) GetAccessRols() []string {
	return []string{RolSuperAdmin, RolAdmin}
}

func (h *UpdateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	updUserId, e := strconv.Atoi(vars["id"])
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	//Get user from database
	var user data.User
	ok, e = h.app.datamodel.ExistUserById(updUserId)
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
	user, e = h.app.datamodel.GetUserById(updUserId)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Reading JSON
	type User struct{ Description, Rol string }
	var update User
	e = ReadJsonFromRequest(r, &update)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Validating description
	if e := data.ValidUserDescription(update.Description); e != nil {
		h.app.elogger.Println(e)
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	//Validating Rol
	if err := h.app.ValidRol(update.Rol); err != nil {
		h.app.elogger.Println(err)
		http.Error(w, "Invalid Rol.", http.StatusBadRequest)
		return
	}

	//SuperAdmin never can be updated
	if user.Rol == RolSuperAdmin {
		h.app.elogger.Println("SuperAdmin users never can be updated")
		http.Error(w, "SuperAdmin users never can be updated.", http.StatusForbidden)
		return
	}
	//SuperAdmin users never can be created
	if update.Rol == RolSuperAdmin {
		h.app.elogger.Println("SuperAdmin users never can be created")
		http.Error(w, "SuperAdmin users never can be created.", http.StatusForbidden)
		return
	}

	//Only SuperAdmin can update admin users
	if user.Rol == RolAdmin && author.Rol != RolSuperAdmin {
		h.app.elogger.Println("Only SuperAdmin can update admin users")
		http.Error(w, "Only SuperAdmin can update admin users.", http.StatusForbidden)
		return
	}
	//Only SuperAdmin can create admin users
	if update.Rol == RolAdmin && author.Rol != RolSuperAdmin {
		h.app.elogger.Println("Only SuperAdmin can create admin users")
		http.Error(w, "Only SuperAdmin can create admin users.", http.StatusForbidden)
		return
	}

	//Update User in database
	user.Rol = update.Rol
	user.Description = update.Description
	e = h.app.datamodel.UpdateUser(&user)

	//Output to JSON
	var output []byte
	if e == nil {
		output, e = json.MarshalIndent(&user, "", "\t\t")
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
