package database

import (
	"log"

	"github.com/google/uuid"
)

func (a *DatabaseAdaptor) Save(data *DummyOrm) (uuid.UUID, error) {
	if err := a.db.Create(data).Error; err != nil {
		log.Println("Cannot create data :%v", err)
		return uuid.Nil, err
	}

	return data.UserID, nil
}

func (a *DatabaseAdaptor) GetByUuid(uuid *uuid.UUID) (DummyOrm, error) {
	var res DummyOrm
	if err := a.db.First(&res, "user_id = ?", uuid).Error; err != nil {
		log.Println("Cannot get data")
		return res, err
	}

	return res, nil
}
