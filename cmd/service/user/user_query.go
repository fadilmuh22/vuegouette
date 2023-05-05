package user

import "restskuy/cmd/db"

// get all user with db
func GetAllUser() ([]User, error) {
	var users []User
	c := db.Connect()

	// retrive user from db using sql query store in users
	result, err := c.Query("SELECT * FROM user")
	if err != nil {
		return nil, err
	}

	defer result.Close()

	for result.Next() {
		var user User
		err := result.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func GetUser(ID string) (User, error) {
	var user User
	c := db.Connect()

	result, err := c.Query("SELECT * FROM user WHERE id = ?", ID)
	if err != nil {
		return user, err
	}

	defer result.Close()

	errorScan := result.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if errorScan != nil {
		return user, errorScan
	}

	return user, nil

}

func CreateUser(user User) (User, error) {
	c := db.Connect()

	_, err := c.Exec("INSERT INTO user (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func UpdateUser(id string, user User) (User, error) {
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
