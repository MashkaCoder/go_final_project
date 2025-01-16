package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	utl     "todo/utils"
	dbfuncs "todo/database"
	schd    "todo/scheduler"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case http.MethodPost:
        id := r.URL.Query().Get("id")
        task, err := utl.CheckExistId(id)
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
        } else {
            idInt, _ := strconv.Atoi(id) // ошибка не обрабатывается, так как была проверка в checkExistId
            err = dbfuncs.DelTaskById(idInt)
            if err != nil{
                http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusMethodNotAllowed)
                return
            }
        }
        err = dbfuncs.UpdateTask(*task)
        if err != nil {
            http.Error(w, `{"error":"Ошибка при обновлении задачи"}`, http.StatusMethodNotAllowed)
            return
        }
        json.NewEncoder(w).Encode(map[string]interface{}{})
    default:
        http.Error(w, `{"error":"Invalid request method"}`, http.StatusMethodNotAllowed)
        return
    }
}