package controller

import (
    "encoding/json"
    "net/http"
    "time"

    "golang-todo-api/model"
    "golang-todo-api/repository"

    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
)

type AuthController struct {
    userRepo repository.UserRepository
}

func NewAuthController(repo repository.UserRepository) *AuthController {
    return &AuthController{userRepo: repo}
}

// Register
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
    var user model.User
    json.NewDecoder(r.Body).Decode(&user)

    // Hash password
    hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }
    user.Password = string(hashed)

    // Simpan ke DB
    err = c.userRepo.Create(user)
    if err != nil {
        http.Error(w, "Error creating user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
}

// Login
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
    var input model.User
    json.NewDecoder(r.Body).Decode(&input)

    user, err := c.userRepo.GetByUsername(input.Username)
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // Cek password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // Buat token JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
    })

    tokenStr, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}
