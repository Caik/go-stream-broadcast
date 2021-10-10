package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RespOk(w http.ResponseWriter, message string, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)

	resp := Response{
		Status:  200,
		Message: message,
		Data:    data,
	}

	body, err := json.Marshal(resp)

	if err != nil {
		w.Write([]byte(`{"status":200}`))
		return
	}

	if _, err = w.Write(body); err != nil {
		log.Printf("Error writing ok to response: %v", err)
	}
}

func RespNotOk(code int, w http.ResponseWriter, message string, data map[string]interface{}) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")

	resp := Response{
		Status:  code,
		Message: message,
		Data:    data,
	}

	body, err := json.Marshal(resp)

	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"status":%d}`, code)))
		return
	}

	if _, err = w.Write(body); err != nil {
		log.Printf("Error writing not ok to response: %v", err)
	}
}
