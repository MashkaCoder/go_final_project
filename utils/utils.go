package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MashkaCoder/go_final_project/constants"
	m "github.com/MashkaCoder/go_final_project/model"
	schd "github.com/MashkaCoder/go_final_project/scheduler"
)

func CheckFormatId(id string) (int, error) {
	if id == "" {
		return 0, fmt.Errorf("не указан идентификатор")
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("идентификатор должен быть числом")
	}
	return idInt, err
}

func ParseHandlerTask(r *http.Request, isPut bool) (*m.Task, error) {
	var task m.Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		return nil, fmt.Errorf("invalid JSON")
	}
	if task.Title == "" {
		return nil, fmt.Errorf("не указан заголовок задачи")

	}
	if task.Date == "" {
		task.Date = time.Now().Format(constants.DateFormat)
	}
	_, err := time.Parse(constants.DateFormat, task.Date)
	if err != nil {
		return nil, fmt.Errorf("дата представлена в неправильном формате")
	}
	if task.Date < time.Now().Format(constants.DateFormat) {
		task.Date = time.Now().Format(constants.DateFormat)
	}
	if task.Repeat != "" {
		task.Date, err = schd.NextDate(time.Now(), task.Date, task.Repeat, false)
		if err != nil {
			return nil, fmt.Errorf("правило повторения указано в неправильном формате")
		}
	}
	if isPut {
		if task.ID == "" {
			return nil, fmt.Errorf("не указан идентификатор задачи для обновления")
		}
		_, err := strconv.Atoi(task.ID)
		if err != nil {
			return nil, fmt.Errorf("идентификатор должен быть числом")
		}
	}
	return &task, nil
}
