// repository/todo_repository.go
package repository

import (
    "database/sql"
    "log"
    "golang-todo-api/model"
    _ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
    var err error
    // Menggunakan database 'postgres'
    connStr := "host=localhost port=5432 user=postgres password=abc123456 dbname=postgres sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
    log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}

type TodoRepository interface {
    GetAll() ([]model.Todo, error)
    GetByID(id int) (model.Todo, error)
    Create(todo model.Todo) (model.Todo, error)
    Update(id int, todo model.Todo) (model.Todo, error)
    Delete(id int) error
}

type todoRepositoryImpl struct {}

func NewTodoRepository() TodoRepository {
    return &todoRepositoryImpl{}
}

func (r *todoRepositoryImpl) GetAll() ([]model.Todo, error) {
    rows, err := db.Query("SELECT id, task, done FROM todos")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var todos []model.Todo
    for rows.Next() {
        var todo model.Todo
        if err := rows.Scan(&todo.ID, &todo.Task, &todo.Done); err != nil {
            return nil, err
        }
        todos = append(todos, todo)
    }

    return todos, nil
}

func (r *todoRepositoryImpl) GetByID(id int) (model.Todo, error) {
    var todo model.Todo
    err := db.QueryRow("SELECT id, task, done FROM todos WHERE id = $1", id).Scan(&todo.ID, &todo.Task, &todo.Done)
    if err != nil {
        return todo, err
    }
    return todo, nil
}

func (r *todoRepositoryImpl) Create(todo model.Todo) (model.Todo, error) {
    err := db.QueryRow("INSERT INTO todos(task, done) VALUES($1, $2) RETURNING id", todo.Task, todo.Done).Scan(&todo.ID)
    if err != nil {
        return todo, err
    }
    return todo, nil
}

func (r *todoRepositoryImpl) Update(id int, todo model.Todo) (model.Todo, error) {
    _, err := db.Exec("UPDATE todos SET task = $1, done = $2 WHERE id = $3", todo.Task, todo.Done, id)
    if err != nil {
        return todo, err
    }
    todo.ID = id
    return todo, nil
}

func (r *todoRepositoryImpl) Delete(id int) error {
    _, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
    return err
}
