package auth

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func CreateNewUser(db *sql.DB, username, password string) error{
	createNewUserQuery := `INSERT INTO users(username, password) VALUES(?,?);`
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(createNewUserQuery, username, passwordHash)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserCredentials(db *sql.DB, uname, passwd string) error{
	fetchPasswdQuery := `SELECT password FROM users WHERE username = ?;`
	var storedPasswd []byte

	passRow := db.QueryRow(fetchPasswdQuery, uname)

	err := passRow.Scan(&storedPasswd)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(storedPasswd, []byte(passwd))

	return err
}

