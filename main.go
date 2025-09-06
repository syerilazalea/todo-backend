package main

import (
    "fmt"
    "log"
    "net/http"

    "golang-todo-api/controller"
    "golang-todo-api/repository"

    "github.com/gorilla/mux"
    "github.com/rs/cors"
	
)

func main() {
    // Inisialisasi koneksi ke database PostgreSQL
    repository.InitDB()

    // Setup repository & controller
    todoRepo := repository.NewTodoRepository()
    todoController := controller.NewTodoController(todoRepo)

    userRepo := repository.NewUserRepository()
    authController := controller.NewAuthController(userRepo)

    // Setup Router
    r := mux.NewRouter()

    // ===== Auth Routes =====
    r.HandleFunc("/register", authController.Register).Methods("POST")
    r.HandleFunc("/login", authController.Login).Methods("POST")

    // ===== Todo Routes (diproteksi JWT) =====
    r.Handle("/todos", controller.AuthMiddleware(http.HandlerFunc(todoController.GetAllTodos))).Methods("GET")
    r.Handle("/todos/{id}", controller.AuthMiddleware(http.HandlerFunc(todoController.GetTodoByID))).Methods("GET")
    r.Handle("/todos", controller.AuthMiddleware(http.HandlerFunc(todoController.CreateTodo))).Methods("POST")
    r.Handle("/todos/{id}", controller.AuthMiddleware(http.HandlerFunc(todoController.UpdateTodo))).Methods("PUT")
    r.Handle("/todos/{id}", controller.AuthMiddleware(http.HandlerFunc(todoController.DeleteTodo))).Methods("DELETE")

    // ===== CORS =====
    handler := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, // asal frontend (Next.js)
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"*"},
        AllowCredentials: true,
    }).Handler(r)

    // Start Server
    fmt.Println("Server berjalan di port 8080...")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
