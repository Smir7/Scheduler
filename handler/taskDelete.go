package handler

import (
	"encoding/json"
	"errors"
	"github.com/smir7/scheduler/constans"
	"github.com/smir7/scheduler/repository"
	"net/http"
)

func TaskDelete(rep repository.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		err := rep.DeleteTask(id)
		if err != nil {
			err := errors.New(" Task with ID not found")
			constans.ErrorResponse.Error = err.Error()
			json.NewEncoder(res).Encode(constans.ErrorResponse)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(map[string]string{}); err != nil {
			http.Error(res, `{"error":"JSON coding error"}`, http.StatusInternalServerError)
			return
		}
	}
}
