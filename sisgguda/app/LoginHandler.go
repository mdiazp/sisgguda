package app

import (
	"net/http"

	"github.com/lamg/regapi"
)

type LoginHandler struct {
	app *App
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var e error

	//Get credentials from request
	credentials, e := regapi.GetCredentials(r)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Authenticate in ldap server
	e = h.app.ldap.Authenticate(credentials.User, credentials.Pass)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, "Bad Credentials", http.StatusUnauthorized)
		return
	}

	//Verificate registration of user
	_, e = h.app.datamodel.GetUserByUsername(credentials.User)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	//Prepare encryption of credentials
	s, e := h.app.crypto.Encrypt(credentials)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Writing Response
	h.app.WriteResponse(w, 200, true, []byte(s))
	return
}
