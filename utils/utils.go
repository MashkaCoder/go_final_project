package utils

import (
    schd "todo/scheduler"
    m    "todo/model"
    "todo/database"
    
    "strconv"
    "fmt"
    "net/http"
    "time"
    "encoding/json"
)

func CheckExistId(id string) (*m.Task, error){
    if id == ""{
        return &m.Task{}, fmt.Errorf("не указан идентификатор")
    }
    idInt, err := strconv.Atoi(id)
    if err != nil{
        return &m.Task{}, fmt.Errorf("идентификатор должен быть числом")
        
    }
    task, err := database.GetTaskById(idInt)
    return &task, err
}

func ParseHandlerTask(r *http.Request, isPut bool) (*m.Task, error){
    var task m.Task
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&task); err != nil {
        return nil, fmt.Errorf("invalid JSON")
    }
    if task.Title == "" {
        return nil, fmt.Errorf("не указан заголовок задачи")
        
    }
    if task.Date == "" {
        task.Date = time.Now().Format("20060102")
    }
    _, err := time.Parse("20060102", task.Date)
    if err != nil {
        return nil, fmt.Errorf("дата представлена в неправильном формате")
    }
    if task.Date < time.Now().Format("20060102") {
        task.Date = time.Now().Format("20060102")
    }
    if task.Repeat != "" {
        task.Date, err = schd.NextDate(time.Now(), task.Date, task.Repeat, false)
        if err != nil {
            return nil, fmt.Errorf("правило повторения указано в неправильном формате")
        }
        }
    if isPut{
        if task.ID == ""{
            return nil, fmt.Errorf("не указан идентификатор задачи для обновления")
        }
        _, err := strconv.Atoi(task.ID)
        if err != nil{
            return nil, fmt.Errorf("идентификатор должен быть числом")
        }
    }
    return &task, nil
}
