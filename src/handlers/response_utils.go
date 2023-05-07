package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func responseJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func badRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func notFound(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(message))
}

func internalError(w http.ResponseWriter, message string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
	fmt.Println(message)
	fmt.Println(err)
}
