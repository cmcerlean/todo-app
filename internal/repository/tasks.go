package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"todo-app/internal/data"

	_ "github.com/go-sql-driver/mysql" // import for driver
)

// DBConnection initalizes a sql.DB instance
func DBConnection(user string, password string, database string) (db *sql.DB, err error) {
	connString := fmt.Sprintf("%s:%s@%s", user, password, database)
	db, err = sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	return
}

// NewTaskDAO returns new instance of TaskDAO
func NewTaskDAO(db *sql.DB) *TaskDAO {
	return &TaskDAO{DB: db}
}

// TaskDAO implements DAO interface to write tasks to DB
type TaskDAO struct {
	DB *sql.DB
}

// InsertTask places a task object into the database
func (dao *TaskDAO) InsertTask(goal string) (id int64, err error) {
	stmt, err := dao.DB.Prepare("INSERT INTO tasks (goal) VALUES (?)")
	if err != nil {
		return
	}

	result, err := stmt.Exec(goal)
	if err != nil {
		return
	}

	return result.LastInsertId()
}

// GetTask finds a task in the database
func (dao *TaskDAO) GetTask(taskID int) (*data.Task, error) {
	task := data.Task{}

	stmt, err := dao.DB.Prepare("SELECT id, goal, completed FROM tasks WHERE id = ?")
	if err != nil {
		return &task, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(taskID).Scan(&task.ID, &task.Goal, &task.Completed)
	if err != nil {
		return &task, data.ErrRecordNotFound
	}

	return &task, nil
}

// GetTasks returns all tasks in database
func (dao *TaskDAO) GetTasks(completed string) (*data.Tasks, error) {

	args := []interface{}{}
	query := "SELECT id, goal, completed FROM tasks"

	if completed != "" {
		query = query + " where completed = ?"

		b, err := strconv.ParseBool(completed)
		if err != nil {
			return &data.Tasks{}, err
		}
		args = append(args, b)
	}

	rows, err := dao.DB.Query(query, args...)
	if err != nil {
		return &data.Tasks{}, err
	}
	defer rows.Close()

	ts := make([]data.Task, 0)
	for rows.Next() {
		var t data.Task
		err := rows.Scan(&t.ID, &t.Goal, &t.Completed)
		if err != nil {
			return &data.Tasks{}, err
		}
		ts = append(ts, t)
	}
	err = rows.Err()
	if err != nil {
		return &data.Tasks{}, err
	}

	return &data.Tasks{Tasks: ts}, nil
}

// DeleteTask removes a task from the database
func (dao *TaskDAO) DeleteTask(taskID int) (err error) {
	stmt, err := dao.DB.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		return
	}

	_, err = stmt.Exec(taskID)
	if err != nil {
		return
	}

	return
}

// UpdateTask updates values of task in database
func (dao *TaskDAO) UpdateTask(taskID int, goal string, completed bool) (rowsUpdated int64, err error) {
	stmt, err := dao.DB.Prepare("UPDATE tasks SET goal = ?, completed = ? WHERE id = ?")
	if err != nil {
		return
	}

	result, err := stmt.Exec(goal, completed, taskID)
	if err != nil {
		return
	}

	return result.RowsAffected()
}
