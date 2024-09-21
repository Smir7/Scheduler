package handler

import (
	"encoding/json"
	"github.com/smir7/scheduler/repository"
	"net/http"
)

func TaskDone(rep repository.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		err := rep.TaskDone(id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(map[string]string{}); err != nil {
			http.Error(res, `{"error":"JSON wrong coding "}`, http.StatusInternalServerError)
			return
		}
	}
}
