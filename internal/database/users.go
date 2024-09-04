package database

import "os"

type User struct {
	Email    string `json:"email"`
	ID       int    `json:"id"`
	Password string `json:"password"`
}

func (db *DB) CreateUser(email, password string) (User, error) {

	if db.EmailExist(email) {
		return User{}, os.ErrExist
	}

	data, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(data.Users) + 1
	user := User{
		ID:       id,
		Email:    email,
		Password: password,
	}

	data.Users[user.ID] = user

	err = db.writeDB(data)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, os.ErrNotExist
	}

	return user, nil
}

func (db *DB) FindUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for key, val := range dbStructure.Users {
		if val.Email == email {
			return dbStructure.Users[key], nil
		}
	}

	return User{}, os.ErrNotExist
}

func (db *DB) EmailExist(email string) bool {

	_, err := db.FindUser(email)
	if err != nil {
		return false
	}

	return true

}

func (db *DB) UpdateUser(id int, email, password string) error {

	data, err := db.loadDB()
	if err != nil {
		return err
	}

	user := data.Users[id]
	user.Email = email
	user.Password = password
	data.Users[id] = user

	err = db.writeDB(data)
	if err != nil {
		return err
	}
	return nil
}
