package helper

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type (
	Meta struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Status  string `json:"status"`
	}

	ResponseFormatter struct {
		Meta Meta        `json:"meta"`
		Data interface{} `json:"data"`
	}
)

func GenerateID() string {
	IDCandidate := uuid.New()

	ID := strings.Replace(IDCandidate.String(), "-", "", -1)
	return ID
}

func JSON(w http.ResponseWriter, data interface{}, status int) {
	dataByte, err := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(status)
	w.Write([]byte(dataByte))
}

func APIResponse(message string, code int, status string, data interface{}) ResponseFormatter {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	response := ResponseFormatter{
		Meta: meta,
		Data: data,
	}

	return response
}
