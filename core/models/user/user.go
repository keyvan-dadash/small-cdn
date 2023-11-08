package user

import (
	"errors"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Uuid     string `pg:",pk"`
	Username string `pg:",unique,notnull"`
	Password string `pg:",notnull"`
}

// CreateUser create user instance based on given username and password
func CreateUser(username string, password string) (*User, error) {
	hashedPassword, err := hashAndSaltPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Uuid:     uuid.NewV4().String(),
		Username: username,
		Password: hashedPassword,
	}, nil
}

// VerifyPassword is function that verfiy given password
func (u *User) VerifyPassword(givenPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(givenPassword))
	return err == nil
}

func hashAndSaltPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", errors.New("cannot generate hash from password")
	}
	return string(hash), nil
}
