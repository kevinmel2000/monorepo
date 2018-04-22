package httpresponse

import (
	"encoding/json"
	"net/http"

	"github.com/lab46/monorepo/gopkg/log"
)

var (
	internalServerErrorResp []byte
)

func init() {
	internalServerErr := map[string]interface{}{
		"errors": []string{"internal server error"},
	}
	internalServerErrorResp, _ = json.Marshal(internalServerErr)
}

// StatusOK response
func StatusOK(w http.ResponseWriter) {
	resp := map[string]interface{}{
		"status": "OK",
	}
	jsonResp, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// InternalServerError response
func InternalServerError(w http.ResponseWriter, errString ...string) {
	status := http.StatusInternalServerError
	resp := map[string]interface{}{
		"errors": errString,
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Failed to marshal message string in Internal Server Error: %s", err.Error())
		jsonResp = internalServerErrorResp
	}
	w.WriteHeader(status)
	w.Write([]byte(jsonResp))
}

// BadRequest response
func BadRequest(w http.ResponseWriter, errString ...string) {
	status := http.StatusBadRequest
	resp := map[string]interface{}{
		"errors": errString,
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Failed to marshal message string in Bad Request : %s", err.Error())
		jsonResp = internalServerErrorResp
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)
	w.Write([]byte(jsonResp))
}

// WithData response
func WithData(w http.ResponseWriter, data interface{}) {
	status := http.StatusOK
	resp := map[string]interface{}{
		"data": data,
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Failed to marshal message string in WithData: %s", err.Error())
		jsonResp = internalServerErrorResp
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)
	w.Write([]byte(jsonResp))
}

// WithObject response
func WithObject(w http.ResponseWriter, object interface{}) {
	status := http.StatusOK
	jsonResp, err := json.Marshal(object)
	if err != nil {
		log.Errorf("Failed to marshal message string in WithObject: %s", err.Error())
		jsonResp = internalServerErrorResp
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)
	w.Write([]byte(jsonResp))
}
