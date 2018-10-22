package app

import "net/http"

type AccessLogHandler struct {
	app  *App
	next http.Handler
}

func (h *AccessLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.app.alogger.Printf("New access from host:  %s\n", r.Host)
	h.next.ServeHTTP(w, r)
}
