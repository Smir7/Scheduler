package handler

import (
	"encoding/json"
	"net/http"

	"github.com/smir7/scheduler/constans"
	"github.com/smir7/scheduler/repository"
)

func TaskPut(rep repository.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var t constans.Task
		err := json.NewDecoder(req.Body).Decode(&t)
		if err != nil {
			http.Error(res, `{"error":"JSON deserialization error"}`, http.StatusBadRequest)
			return
		}
		err = rep.UpdateTask(t)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(map[string]string{}); err != nil {
			http.Error(res, `{"error":"JSON coding error"}`, http.StatusInternalServerError)
			return
		}
	}
}
