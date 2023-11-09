package auth

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/sod-lol/small-cdn/core/models/user"
)

func GenerateUser(signUpJson signUpJsonExpect) (int, error) {
	tempUser, err := user.CreateUser(signUpJson.Username, signUpJson.Password)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := user.UserRepository.InsertUser(tempUser); err != nil {

		if !errors.Is(err, user.ErrDublicateUser) {
			logrus.Errorf("Cannot insert user. error: %v", err)
			return http.StatusInternalServerError, errors.New("cannot signup")
		}

		return http.StatusBadRequest, errors.New("user with given username already exists")
	}

	return http.StatusCreated, nil
}
