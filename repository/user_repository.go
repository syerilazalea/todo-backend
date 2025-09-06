package repository

import (
    "golang-todo-api/model"
)

type UserRepository interface {
    Create(user model.User) error
    GetByUsername(username string) (model.User, error)
}

type userRepositoryImpl struct{}

func NewUserRepository() UserRepository {
    return &userRepositoryImpl{}
}

func (r *userRepositoryImpl) Create(user model.User) error {
    _, err := db.Exec("INSERT INTO users(username, password) VALUES($1, $2)", user.Username, user.Password)
    return err
}

func (r *userRepositoryImpl) GetByUsername(username string) (model.User, error) {
    var user model.User
    err := db.QueryRow("SELECT id, username, password FROM users WHERE username=$1", username).
        Scan(&user.ID, &user.Username, &user.Password)
    return user, err
}
