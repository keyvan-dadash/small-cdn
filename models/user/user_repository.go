package user

import "github.com/go-pg/pg/v10"

type UserRepoInterface interface {
	InsertUser(*User) error
	RetrieveUser(*User) error
	UpdateUser(*User) error
	DeleteUser(string) error
}

var UserRepository *UserRepo

type UserRepo struct {
	*pg.DB
}

func CreateUserRepo(db *pg.DB) {
	UserRepository = &UserRepo{
		DB: db,
	}
}

func (ur *UserRepo) InsertUser(user *User) error {
	_, err := ur.DB.Model(user).Insert()
	return err
}

func (ur *UserRepo) RetrieveUser(user *User) error {
	err := ur.DB.Model(user).Where("username = ?", user.Username).Select()
	return err
}

func (ur *UserRepo) UpdateUser(user *User) error {
	_, err := ur.DB.Model(user).WherePK().Update()
	return err
}

func (ur *UserRepo) DeleteUser(username string) error {
	_, err := ur.DB.Model(((*User)(nil))).Where("username = ?", username).Delete()
	return err
}
