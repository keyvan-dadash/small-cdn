package user

import (
	"gorm.io/gorm"

	"github.com/sod-lol/small-cdn/core/models/cache"
)

type UserRepoInterface interface {
	InsertUser(*User) error
	RetrieveUser(*User) error
	RetrieveUserPreloadCacheLogs(string) (error, []cache.CacheLog)
	UpdateUser(*User) error
	DeleteUser(string) error
}

var UserRepository *UserRepo

type UserRepo struct {
	*gorm.DB
}

func CreateUserRepo(db *gorm.DB) {
	UserRepository = &UserRepo{
		DB: db,
	}
}

func (ur *UserRepo) InsertUser(user *User) error {
	result := ur.DB.Create(user)
	return result.Error
}

func (ur *UserRepo) RetrieveUser(user *User) error {
	result := ur.DB.Where("username = ?", user.Username).First(user)
	return result.Error
}

func (ur *UserRepo) RetrieveUserPreloadCacheLogs(username string) (err error, logs []cache.CacheLog) {
	var user User
	result := ur.DB.Model(&User{}).Preload("CacheLogs").Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, []cache.CacheLog{}
	}
	return result.Error, user.CacheLogs
}

func (ur *UserRepo) UpdateUser(user *User) error {
	result := ur.DB.Save(user)
	return result.Error
}

func (ur *UserRepo) DeleteUser(username string) error {
	result := ur.DB.Where("username = ?", username).Delete(&User{})
	return result.Error
}
