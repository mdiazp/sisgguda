package app

import (
	"context"
	"net/http"
)

type AuthHandler struct {
	app  *App
	next http.Handler
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Verificate header authHd
	username, e := h.app.crypto.Decrypt(r)
	if e != nil {
		h.app.elogger.Println("Error crypto.Decrypt")
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	//Look into ldap server the username
	//In case that user by invalid in ldap server after login in sisgguda
	if e := h.app.ValidAdUser(username); e != nil {
		h.app.elogger.Println("Undefined user in ldap server.")
		http.Error(w, "Invalid user in ldap server", http.StatusForbidden)
		return
	}

	//Get perfil of user from database
	author, e := h.app.datamodel.GetUserByUsername(username)
	if e != nil {
		h.app.elogger.Println(e)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Put user in context for other handlers
	ctx := context.WithValue(r.Context(), "author", author)

	//next handler
	h.next.ServeHTTP(w, r.WithContext(ctx))
}

func MustAuth(app *App, handler http.Handler) http.Handler {
	return &AuthHandler{app: app, next: handler}
}
