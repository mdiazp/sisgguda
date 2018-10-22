package app

import (
	"encoding/json"
	"errors"
	"net/http"
)

// version
func version() string {
	return "0.1"
}

func containString(slice []string, s string) bool {
	for _, x := range slice {
		if x == s {
			return true
		}
	}
	return false
}

func ReadJsonFromRequest(r *http.Request, objs ...interface{}) (e error) {
	decoder := json.NewDecoder(r.Body)
	for i, obj := range objs {
		e = decoder.Decode(obj)
		if e != nil {
			e = errors.New("Error decoding element " + string(i) + "   error: " + e.Error())
			return
		}
	}
	return
}
