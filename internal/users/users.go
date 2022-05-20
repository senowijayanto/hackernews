package users

import (
	"database/sql"
	"log"

	db "github.com/senowijayanto/hackernews/internal/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

type WrongUsernameOrPasswordError struct{}

func (user *User) Create() {
	stmt, err := db.Db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	print(stmt)
	if err != nil {
		log.Fatal(err)
	}
	hashPassword, err := HashPassword(user.Password)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(user.Username, hashPassword)
	if err != nil {
		log.Fatal(err)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserIdByUsername(username string) (int, error) {
	stmt, err := db.Db.Prepare("SELECT id from users WHERE username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(username)

	var id int
	err = row.Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return id, nil
}

func (user *User) Authenticate() bool {
	stmt, err := db.Db.Prepare("SELECT password from users WHERE username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(user.Username)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	return CheckPasswordHash(user.Password, hashedPassword)
}

func (m *WrongUsernameOrPasswordError) Error() string {
	return "Wrong username or password"
}
