package database

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseAdaptor struct {
	db *gorm.DB
}

func NewDatabaseAdaptor(conn *sql.DB) (*DatabaseAdaptor, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("Can't connect database (gorm) : %v", err)
	}

	return &DatabaseAdaptor{
		db: db,
	}, nil
}
