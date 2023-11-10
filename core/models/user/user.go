package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/sod-lol/small-cdn/core/models/cache"
)

type User struct {
	gorm.Model
	Username  string           `gorm:"unique"`
	Password  string           `gorm:"not null"`
	CacheLogs []cache.CacheLog `gorm:"foreignKey:UserID"`
}

// CreateUser create user instance based on given username and password
func CreateUser(username string, password string) (*User, error) {
	hashedPassword, err := hashAndSaltPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:  username,
		Password:  hashedPassword,
		CacheLogs: []cache.CacheLog{},
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
