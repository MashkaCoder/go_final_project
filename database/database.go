package database

import (
    "database/sql"
    "errors"
    "fmt"
    "log"
    "os"
    "path/filepath"
	m "todo/model"

    _ "modernc.org/sqlite"
)

var db *sql.DB

func GetDB() *sql.DB {
    return db
}

func CreateTable() error {
    appPath, err := os.Getwd()
    if err != nil {
        return err
    }
    dbFile := filepath.Join(appPath, "scheduler.db")
    _, err = os.Stat(dbFile)

    if os.IsNotExist(err) {
        log.Println("Создаём новую базу данных и таблицу...")
        db, err = sql.Open("sqlite", dbFile)
        if err != nil {
            return err
        }
        createTableSQL := `
            CREATE TABLE IF NOT EXISTS scheduler (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                date TEXT NOT NULL,
                title TEXT NOT NULL,
                comment TEXT,
                repeat TEXT CHECK(length(repeat) <= 128)
            );
        `
        _, err = db.Exec(createTableSQL)
        if err != nil {
            return err
        }
    } else if err != nil {
        return err
    } else {
        db, err = sql.Open("sqlite", dbFile)
        if err != nil {
            return err
        }
        err = db.Ping()
        if err != nil {
            return fmt.Errorf("не удалось подключиться к базе данных: %v", err)
        }
    }
    return nil
}


func GetTasks() ([]m.Task, error) {
    if db == nil {
        return nil, errors.New("не установлене соедениние с базой данных")
    }
    query := `SELECT * FROM scheduler ORDER BY date ASC LIMIT 50`
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    var tasks []m.Task
    for rows.Next(){
        var task m.Task
        if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
            return nil, err
        }
        tasks = append(tasks, task)
    }
    if err := rows.Err(); err != nil {
		return nil, err
	}
    if len(tasks) == 0 {
        return []m.Task{}, nil
    }
    return tasks, nil
}

func GetTaskById(id int) (m.Task, error){
    if db == nil {
        return m.Task{}, errors.New("не установлене соедениние с базой данных")
    }
    query := `SELECT * FROM scheduler 
              WHERE id = ?`
    var task m.Task
    err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
    if err == sql.ErrNoRows {
        return m.Task{}, fmt.Errorf("задача с ID %d не найдена", id)
    } 
    if err != nil {
        return m.Task{}, err
    }
    return task, nil
}

func InsertTask(task m.Task) (int64, error) {
    if db == nil {
        return 0, errors.New("не установлено соединение с базой данных")
    }
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateTask(task m.Task) error {
    if db == nil {
        return errors.New("не установлене соедениние с базой данных")
    }
    query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
    _, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
    return err
}

func DelTaskById(id int) error{
    if db == nil {
        return errors.New("не установлене соедениние с базой данных")
    }
    query := `DELETE from scheduler WHERE id = ?`
    _, err := db.Exec(query, id)
    return err
}