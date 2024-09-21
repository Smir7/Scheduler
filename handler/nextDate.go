package handler

import (
	"net/http"
	"time"

	"github.com/smir7/scheduler/constans"
	"github.com/smir7/scheduler/task"
)

func NextDate(res http.ResponseWriter, req *http.Request) {
	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	nowTime, err := time.Parse(constans.DateFormat, now)
	if err != nil {
		http.Error(res, "format date error", http.StatusBadRequest)
		return
	}
	nextDate, err := task.NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = res.Write([]byte(nextDate))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

}
