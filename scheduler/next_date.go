package scheduler

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string, isPost bool) (string, error) {
    if repeat == "" {
        return date, nil
    }
    if !(strings.HasPrefix(repeat, "d") || strings.HasPrefix(repeat, "y")) {
        return "", errors.New("не поддерживаемый формат")
    }
    dateStart, err := time.Parse("20060102", date)
    if err != nil {
        return "", err
    }
    res := dateStart
    if repeat == "y" {
        years := now.Year() - dateStart.Year()
        if years <= 0 {
            years = 1 // Если текущий год меньше года начала, через 1 год
        }
        res = res.AddDate(years, 0, 0)
    } else if strings.HasPrefix(repeat, "d")  {
        days := strings.Split(repeat, " ")
        if len(days) != 2 {
            return "", errors.New("не указан интервал в днях")
        }
        day, err := strconv.Atoi(days[1])
        if err != nil {
            return "", err
        }
        if day > 400 {
            return "", errors.New("превышен максимально допустимый интервал в 400 дней")
        }
        if isPost {
            res = res.AddDate(0, 0, day)
            return res.Format("20060102"), nil
        }
        if res.Format("20060201") == now.Format("20060201"){ ////
            return now.Format("20060102"), nil
        } else if res.Before(now) {
            for res.Before(now){
                res = res.AddDate(0, 0, day)
            }
        } else { 
            res = res.AddDate(0, 0, day)
        }
    }
    return res.Format("20060102"), nil
}