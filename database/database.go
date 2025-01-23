package database

import (
	"database/sql"
	"errors"
	"fmt"
	// "log"
	"os"
	"path/filepath"

	"github.com/MashkaCoder/go_final_project/constants"
	"github.com/MashkaCoder/go_final_project/model"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func createTable(conn *sql.DB) error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT CHECK(length(repeat) <= 128)
		);
		CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	`
	_, err := conn.Exec(createTableSQL)
	return err
}

func NewDB() (*DB, error) {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		appPath, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		dbFile = filepath.Join(appPath, "scheduler.db")
	}
	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}

	err = createTable(conn)
	if err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

func (db *DB) GetTasks() ([]model.Task, error) {
	if db.conn == nil {
		return nil, errors.New("не установлено соединение с базой данных")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
	rows, err := db.conn.Query(query, constants.TaskLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return []model.Task{}, nil
	}
	return tasks, nil
}

func (db *DB) GetTaskById(id int) (model.Task, error) {
	if db.conn == nil {
		return model.Task{}, errors.New("не установлено соединение с базой данных")
	}

	query := `SELECT * FROM scheduler WHERE id = ?`
	var task model.Task
	err := db.conn.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		return model.Task{}, fmt.Errorf("задача с ID %d не найдена", id)
	}
	if err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (db *DB) InsertTask(task model.Task) (int64, error) {
	if db.conn == nil {
		return 0, errors.New("не установлено соединение с базой данных")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.conn.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db *DB) UpdateTask(task model.Task) error {
	if db.conn == nil {
		return errors.New("не установлено соединение с базой данных")
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.conn.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowAff, err := res.RowsAffected()
	if rowAff != 1 {
		return errors.New("задача не найдена")
	}
	return nil
}

func (db *DB) DelTaskById(id int) error {
	if db.conn == nil {
		return errors.New("не установлено соединение с базой данных")
	}

	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.conn.Exec(query, id)
	if err != nil {
		return err
	}

	rowAff, err := res.RowsAffected()
	if rowAff != 1 {
		return errors.New("задача не найдена")
	}
	return nil
}

func (db *DB) GetConnection() *sql.DB {
	return db.conn
}
