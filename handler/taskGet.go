package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/smir7/scheduler/constans"
	"github.com/smir7/scheduler/repository"
)

func TaskGet(rep repository.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		id := req.URL.Query().Get("id")
		task, err := rep.GetTask(id)
		if err != nil {
			err := errors.New("task with ID not found")
			constans.ErrorResponse.Error = err.Error()
			json.NewEncoder(res).Encode(constans.ErrorResponse)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(task); err != nil {
			http.Error(res, `{"error":"JSON coding error"}`, http.StatusInternalServerError)
			return
		}
	}
}
