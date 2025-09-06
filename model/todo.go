// model/todo.go
package model

type Todo struct {
    ID   int    `json:"id"`
    Task string `json:"task"`
    Done bool   `json:"done"`
}
