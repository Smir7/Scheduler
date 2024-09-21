package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/smir7/scheduler/constans"
	"github.com/smir7/scheduler/repository"
)

func TasksGet(rep repository.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		searchParams := req.URL.Query().Get("search")
		tasks, err := rep.GetTasks(searchParams)
		if err != nil {
			if err != nil {
				err := errors.New("tasks getting error")
				constans.ErrorResponse.Error = err.Error()
				json.NewEncoder(res).Encode(constans.ErrorResponse)
				return
			}
		}
		response := map[string][]constans.Task{
			"tasks": tasks,
		}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, `{"error":"JSON coding error"}`, http.StatusInternalServerError)
			return
		}
	}
}
