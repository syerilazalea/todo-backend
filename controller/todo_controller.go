package controller

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "golang-todo-api/repository"
    "golang-todo-api/model"
)

type TodoController struct {
    todoRepo repository.TodoRepository
}

func NewTodoController(repo repository.TodoRepository) *TodoController {
    return &TodoController{todoRepo: repo}
}

// GET /todos
func (c *TodoController) GetAllTodos(w http.ResponseWriter, r *http.Request) {
    todos, err := c.todoRepo.GetAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(todos)
}

// GET /todos/{id}
func (c *TodoController) GetTodoByID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    todo, err := c.todoRepo.GetByID(id)
    if err != nil {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(todo)
}

// POST /todos
func (c *TodoController) CreateTodo(w http.ResponseWriter, r *http.Request) {
    var newTodo model.Todo
    if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    createdTodo, err := c.todoRepo.Create(newTodo)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdTodo)
}

// PUT /todos/{id}
func (c *TodoController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var updatedTodo model.Todo
    if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    todo, err := c.todoRepo.Update(id, updatedTodo)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(todo)
}

// DELETE /todos/{id}
func (c *TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    if err := c.todoRepo.Delete(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
