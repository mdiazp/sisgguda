package app

import (
	"encoding/json"
	"net/http"
	"time"

	"gitlab.com/manuel.diaz/sisgguda/data"
)

type AddUserHandler struct {
	app *App
}

func (h *AddUserHandler) GetAccessRols() []string {
	return []string{RolSuperAdmin, RolAdmin}
}

func (h *AddUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	type User struct{ Username, Description, Rol string }
	var newUser User
	e = ReadJsonFromRequest(r, &newUser)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Validating user in ldap server
	if e := h.app.ValidAdUser(newUser.Username); e != nil {
		h.app.elogger.Println(e)
		http.Error(w, "Invalid user in ldap server.", http.StatusBadRequest)
		return
	}

	//Validating not conflict with equal usernames
	ok, e = h.app.datamodel.ExistUserByUsername(newUser.Username)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if ok {
		h.app.elogger.Println(e)
		http.Error(w, "Username already exists.", http.StatusConflict)
		return
	}

	//Validating description
	if e := data.ValidUserDescription(newUser.Description); e != nil {
		h.app.elogger.Println(e)
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	//Validating Rol
	if e := h.app.ValidRol(newUser.Rol); e != nil {
		h.app.elogger.Println(e)
		http.Error(w, "Invalid Rol.", http.StatusBadRequest)
		return
	}

	//SuperAdmin users never can be created
	if newUser.Rol == RolSuperAdmin {
		h.app.elogger.Println("SuperAdmin users never can be created")
		http.Error(w, "SuperAdmin users never can be created.", http.StatusForbidden)
		return
	}

	//Only SuperAdmin can insert new admin users
	if newUser.Rol == RolAdmin && author.Rol != RolSuperAdmin {
		h.app.elogger.Println("Only SuperAdmin can insert new admin users")
		http.Error(w, "Only SuperAdmin can insert new admin users.", http.StatusForbidden)
		return
	}

	//Add User to database
	user := data.User{Username: newUser.Username, Description: newUser.Description, Rol: newUser.Rol, CreatedAt: time.Now()}
	e = h.app.datamodel.AddUser(&user)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Output to JSON
	output, e := json.MarshalIndent(&user, "", "\t\t")
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Writing Response
	h.app.WriteResponse(w, http.StatusCreated, true, output)
	return
}
