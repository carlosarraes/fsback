package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := ErrorResponse{Message: message}
	jsonErr, _ := json.Marshal(err)
	i, _ := w.Write(jsonErr)
	fmt.Println(i)
}
