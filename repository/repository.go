package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"

	"github.com/smir7/scheduler/constans"
	"github.com/smir7/scheduler/task"
)

type Repository struct {
	db *sql.DB
}

func CreateDatabase() *sql.DB {
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := os.Getenv("TODO_DBFILE")
	dbFile = filepath.Join(filepath.Dir(path), dbFile)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Loading repository from", "dbFile")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install {
		_, err = db.ExecContext(
			context.Background(), `CREATE TABLE IF NOT EXISTS scheduler 
			    (id INTEGER PRIMARY KEY AUTOINCREMENT, 
			    date CHAR(8) NOT NULL, 
				title VARCHAR(128) NOT NULL, 
				comment VARCHAR(256) NOT NULL, 
				repeat VARCHAR(128) NOT NULL )`)

		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(`CREATE INDEX
				IF NOT EXISTS sheduler_date
				ON scheduler (date);`)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Database installed")
	} else {
		log.Println("Database already exists")
	}
	return db
}

func NewDatabase(db *sql.DB) Repository {
	return Repository{db: db}
}

// Add task
func (rep *Repository) AddTask(t constans.Task) (string, error) {
	var err error
	if t.Title == "" {
		return "", fmt.Errorf(`{"error":"do not title"}`)
	}

	if t.Date == "" {
		t.Date = time.Now().Format(constans.DateFormat)
	}

	_, err = time.Parse(constans.DateFormat, t.Date)
	if err != nil {
		return "", fmt.Errorf(`{"error":"wrong date format"}`)
	}

	if t.Date < time.Now().Format(constans.DateFormat) {
		if t.Repeat != "" {
			nextDate, err := task.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				return "", fmt.Errorf(`{"error":"wrong repeat date"}`)
			}
			t.Date = nextDate
		} else {
			t.Date = time.Now().Format(constans.DateFormat)
		}
	}

	// Add Task in Database
	query := `INSERT INTO scheduler (date, title,comment, repeat) VALUES (?,?,?,?)`
	result, err := rep.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return "", fmt.Errorf(`{"error":"error adding task "}`)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf(`{"error":"error getting last inserted id"}`)
	}
	return fmt.Sprintf("%d", id), nil
}

func (rep *Repository) GetTasks(search string) ([]constans.Task, error) {
	var t constans.Task
	var tasks []constans.Task
	var rows *sql.Rows
	var err error
	if search == "" {
		rows, err = rep.db.Query("SELECT id, date, title, comment, repeat FROM scheduler  ORDER BY date LIMIT ?", constans.CountTasks)
	} else if date, error := time.Parse(constans.DateFormat, search); error == nil {

		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
		rows, err = rep.db.Query(query, date.Format(constans.DateFormat), constans.CountTasks)
	} else {
		search = "%%%" + search + "%%%"
		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		rows, err = rep.db.Query(query, date.Format(constans.DateFormat), constans.CountTasks)
	}
	if err != nil {
		return []constans.Task{}, fmt.Errorf(`{"error":"wrong request"}`)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err = rows.Err(); err != nil {
			return []constans.Task{}, fmt.Errorf(`{"error":"error scanning row"}`)
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return []constans.Task{}, fmt.Errorf(`{"error":"error scanning row"}`)
	}
	if len(tasks) == 0 {
		tasks = []constans.Task{}
	}

	return tasks, nil
}

func (rep *Repository) GetTask(id string) (constans.Task, error) {
	var t constans.Task

	if id == "" {
		return constans.Task{}, fmt.Errorf(`{"error":"id not found"}`)
	}
	row := rep.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return constans.Task{}, fmt.Errorf(`{"error":"task not found"}`)
	}
	return t, nil
}

func (rep *Repository) UpdateTask(t constans.Task) error {
	if t.ID == "" {
		return fmt.Errorf(`{"error":"wrong ID"}`)
	}

	if t.Title == "" {
		return fmt.Errorf(`{"error":"wrong title"}`)
	}

	if t.Date == "" {
		t.Date = time.Now().Format(constans.DateFormat)
	}

	_, err := time.Parse(constans.DateFormat, t.Date)
	if err != nil {
		return fmt.Errorf(`{"error":"Wrong date format"}`)
	}

	if t.Date < time.Now().Format(constans.DateFormat) {
		if t.Repeat != "" {
			nextDate, err := task.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				return fmt.Errorf(`{"error":"Wrong repeat date"}`)
			}
			t.Date = nextDate
		} else {
			t.Date = time.Now().Format(constans.DateFormat)
		}
	}

	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	result, err := rep.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)

	if err != nil {
		return fmt.Errorf(`{"error":"Task is not found"}`)
	}

	rowsChange, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(`{"error":"changes row cannot be counted"}`)
	}

	if rowsChange == 0 {
		return fmt.Errorf(`{"error":"task is not found"}`)
	}
	return nil
}

func (rep *Repository) TaskDone(id string) error {
	var t constans.Task

	t, err := rep.GetTask(id)
	if err != nil {
		return err
	}

	if t.Repeat == "" {
		err := rep.DeleteTask(id)
		if err != nil {
			return err
		}

	} else {
		next, err := task.NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			return err
		}

		t.Date = next
		err = rep.UpdateTask(t)
		if err != nil {
			return err
		}
	}
	return nil
}

// delete task
func (rep *Repository) DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf(`{"error":"wrong ID"}`)
	}

	query := "DELETE FROM scheduler WHERE id == ?"
	result, err := rep.db.Exec(query, id)

	if err != nil {
		return fmt.Errorf(`{"error":"Wrong delete task"}`)
	}

	rowsChange, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf(`{"error":"Error changes row"}`)
	}

	if rowsChange == 0 {
		return fmt.Errorf(`{"error":"task is not found"}`)
	}
	return nil
}
