package helper

import (
	"encoding/json"
	"net/http"
)

func ResponseWithError(w http.ResponseWriter, code int, msg string) {
	Response(w, code, map[string]string{"error": msg})
}

func Response(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
