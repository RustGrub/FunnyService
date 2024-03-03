package utils

import "net/http"

/*
func Write200(body []byte, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}
*/

func Write400(err string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(err))
}

func Write500(err string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(err))
}
