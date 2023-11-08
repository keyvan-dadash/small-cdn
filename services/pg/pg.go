package postgres

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/sirupsen/logrus"
)

type DictoDBInfo struct {
	Addr     string
	Username string
	Password string
	DBName   string
	Port     int32
}

func CreateDictoDBClient(dbInfo *DictoDBInfo) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", dbInfo.Addr, dbInfo.Port),
		User:     dbInfo.Username,
		Password: dbInfo.Password,
		Database: dbInfo.DBName,
	})

	ctx := context.Background()

	_, err := db.ExecContext(ctx, "SELECT 1")
	if err != nil {
		logrus.Fatalf("[Falat] Cannot connect to database. err = %v", err)
		panic(err)
	}

	return db
}

func CreateSchema(db *pg.DB, models ...interface{}) error {
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
