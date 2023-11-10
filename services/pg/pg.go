package postgres

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CDNDBInfo struct {
	Addr     string
	Username string
	Password string
	DBName   string
	Port     int32
	Log      bool
}

func (c *CDNDBInfo) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Tehran",
		c.Addr, c.Username, c.Password, c.DBName, c.Port)
}

func CreateCDNDBClient(dbInfo *CDNDBInfo) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbInfo.GetDSN()), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("[Fatal] Cannot connect to database. err = %v", err)
		panic(err)
	}

	if dbInfo.Log {
		db.Logger.LogMode(logger.Info)
	} else {
		db.Logger.LogMode(logger.Error)
	}

	return db
}

func CreateSchema(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models)
}
