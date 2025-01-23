package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MashkaCoder/go_final_project/constants"
	schd "github.com/MashkaCoder/go_final_project/scheduler"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, "переданы не все параметры", http.StatusBadRequest)
		return
	}
	now, err := time.Parse(constants.DateFormat, nowStr)
	if err != nil {
		http.Error(w, "не верный формат даты", http.StatusBadRequest)
		return
	}
	nextDate, err := schd.NextDate(now, dateStr, repeat, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, nextDate)
}
