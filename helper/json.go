package helper

import (
	"encoding/json"
	"net/http"

	"github.com/mozartmuhammad/julo-be-test/model/web"
)

func WriteSuccess(w http.ResponseWriter, data interface{}) {
	raw, _ := json.Marshal(web.WebResponse{
		Status: "success",
		Data:   data,
	})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(raw)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, err interface{}) {
	raw, _ := json.Marshal(web.WebResponse{
		Status: "fail",
		Data: map[string]interface{}{
			"error": err,
		},
	})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(raw)
}
