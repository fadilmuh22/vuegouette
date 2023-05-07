package user

import (
	"log"

	"github.com/fadilmuh22/restskuy/cmd/db"
	"github.com/fadilmuh22/restskuy/cmd/model"
)

// get all user with db
func GetAllUser() ([]model.User, error) {
	var users []model.User
	c := db.Connect()

	// retrive user from db using sql query store in users
	result, err := c.Query("SELECT id, name, email, password FROM user")
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

func GetUser(ID string) (model.User, error) {
	var user model.User
	c := db.Connect()

	result, err := c.Query("SELECT id, name, email, password FROM user WHERE id = ?", ID)
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

func CreateUser(user model.User) (model.User, error) {
	c := db.Connect()

	result, err := c.Exec("INSERT INTO user (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
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

func UpdateUser(id string, user model.User) (model.User, error) {
	c := db.Connect()

	_, err := c.Exec("UPDATE user SET name = ?, email = ?, password = ? WHERE id = ?", user.Name, user.Email, user.Password, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func DeleteUser(id string) error {
	c := db.Connect()

	_, err := c.Exec("DELETE FROM user WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
