package helper

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
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
