package service

import (
	"database/sql"
	"log"

	"github.com/fadilmuh22/restskuy/cmd/model"
)

type UserService struct {
	Con *sql.DB
}

// get all user with db
func (s UserService) GetAllUser() ([]model.User, error) {
	var users []model.User

	// retrive user from db using sql query store in users
	result, err := s.Con.Query("SELECT id, name, email, password FROM user")
	if err != nil {
		return nil, err
	}

	defer result.Close()

	for result.Next() {
		var user model.User
		err := result.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s UserService) GetUser(ID string) (model.User, error) {
	var user model.User

	result, err := s.Con.Query("SELECT id, name, email, password FROM user WHERE id = ?", ID)
	if err != nil {
		return user, err
	}

	defer result.Close()

	for result.Next() {
		err := result.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			log.Fatal(err)
			return user, err
		}
	}

	return user, nil

}

func (s UserService) CreateUser(user model.User) (model.User, error) {
	result, err := s.Con.Exec("INSERT INTO user (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
	if err != nil {
		return user, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return user, err
	}

	user.ID = int(id)

	return user, nil
}

func (s UserService) UpdateUser(id string, user model.User) (model.User, error) {
	_, err := s.Con.Exec("UPDATE user SET name = ?, email = ?, password = ? WHERE id = ?", user.Name, user.Email, user.Password, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s UserService) DeleteUser(id string) error {
	_, err := s.Con.Exec("DELETE FROM user WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
