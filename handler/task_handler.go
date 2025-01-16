package handler

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    dbfuncs "todo/database"
    utl     "todo/utils"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case http.MethodPost:
        task, err := utl.ParseHandlerTask(r, false)
        if err != nil{
            http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
            return
        }
        id, err := dbfuncs.InsertTask(*task)
        if err != nil {
            http.Error(w, `{"error":"Ошибка при добавлении задачи в базу данных"}`, http.StatusInternalServerError)
            return
        }
        resp := map[string]interface{}{
            "id": id,
        }
        json.NewEncoder(w).Encode(resp)
    case http.MethodGet:
        id  := r.URL.Query().Get("id")
        task, err := utl.CheckExistId(id)
        if err != nil{
            http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusMethodNotAllowed)
            return
        }
        if err := json.NewEncoder(w).Encode(*task); err != nil {
            http.Error(w, `{"error": "Ошибка при отправке ответа"}`, http.StatusInternalServerError)
        }
    case http.MethodPut:
        task, err := utl.ParseHandlerTask(r, true)
        if err != nil{
            http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
            return
        }
        _, err = utl.CheckExistId(task.ID)
        if err != nil{
            http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusMethodNotAllowed)
            return
        }
        err = dbfuncs.UpdateTask(*task)
        if err != nil{
            http.Error(w, `{"error":"Задача не найдена"}`, http.StatusInternalServerError)
        }
        json.NewEncoder(w).Encode(map[string]interface{}{})
    case http.MethodDelete:
        id  := r.URL.Query().Get("id")
        _, err := utl.CheckExistId(id)
        if err != nil{
            http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusMethodNotAllowed)
            return
        }
        idInt, _ := strconv.Atoi(id) // не обрабатываю ошибку, тк уже выполнялась проверка в checkExistId()
        err = dbfuncs.DelTaskById(idInt)
        if err != nil{
            http.Error(w, `{"error":"Ошибка при удалении задач"}`, http.StatusMethodNotAllowed)
            return
        }
        json.NewEncoder(w).Encode(map[string]interface{}{})
    default:
        http.Error(w, `{"error":"Invalid request method"}`, http.StatusMethodNotAllowed)
        return
    }
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch r.Method{
    case http.MethodGet:
        tasks, err := dbfuncs.GetTasks()
        if err != nil{
            http.Error(w, `{"error":"Ошибка при получении задач"}`, http.StatusMethodNotAllowed)
            return
        }
        response := map[string]interface{}{
            "tasks": tasks,
        }
        json.NewEncoder(w).Encode(response)
    default:
        http.Error(w, `{"error":"Invalid request method"}`, http.StatusMethodNotAllowed)
        return
    }    
}