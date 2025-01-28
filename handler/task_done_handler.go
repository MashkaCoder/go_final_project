package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	schd "github.com/MashkaCoder/go_final_project/scheduler"
	utl "github.com/MashkaCoder/go_final_project/utils"
)

func (h *Handler) TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPost:
		id := r.URL.Query().Get("id")
		idInt, err := utl.CheckFormatId(id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusMethodNotAllowed)
			return
		}
		task, err := h.DB.GetTaskById(idInt)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusMethodNotAllowed)
			return
		}

		if task.Repeat != "" {
			dateNew, err := schd.NextDate(time.Now(), task.Date, task.Repeat, true)
			if err != nil {
				http.Error(w, `{"error":"Ошибка при вычислении следующей даты"}`, http.StatusMethodNotAllowed)
				return
			}
			task.Date = dateNew
			err = h.DB.UpdateTask(task)
			if err != nil {
				http.Error(w, `{"error":"Ошибка при обновлении задачи"}`, http.StatusMethodNotAllowed)
				return
			}
		} else {
			idInt, _ := strconv.Atoi(id) // ошибка не обрабатывается, так как была проверка в checkExistId
			err = h.DB.DelTaskById(idInt)
			if err != nil {
				http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusMethodNotAllowed)
				return
			}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
	default:
		http.Error(w, `{"error":"Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}
}
