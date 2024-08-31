package testutil

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetTestDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	sqlDB, sqlMock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open(mysql.New(
		mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		},
	), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return db, sqlMock, nil
}
