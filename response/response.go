package response

import (
	"encoding/json"
	models "go-jwt-api/model"
	"net/http"
)

func ResponseErr(w http.ResponseWriter, statusCode int) {
	jData, err := json.Marshal(models.Error{
		Status:  statusCode,
		Message: http.StatusText(statusCode),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)

}

func ResponseOk(w http.ResponseWriter, data interface{}) {
	if data == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}
