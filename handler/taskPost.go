package handler

import (
	"encoding/json"
	"net/http"

	"github.com/smir7/scheduler/constans"
	"github.com/smir7/scheduler/repository"
)

func TaskPost(rep repository.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var t constans.Task
		err := json.NewDecoder(req.Body).Decode(&t)
		if err != nil {
			http.Error(res, `{"error":"JSON deserialization error"}`, http.StatusBadRequest)
			return
		}
		id, err := rep.AddTask(t)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		response := constans.Response{ID: id}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, `{"error":"JSON coding error"}`, http.StatusInternalServerError)
			return
		}
	}
}
