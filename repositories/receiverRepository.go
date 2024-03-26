package repositories

import (
	"fmt"
	"receiver/handlers"
	"receiver/models"
)

type ReceiverRepository struct {
	Handler handlers.Handler
}

func New(handler handlers.Handler) ReceiverRepository {
	return ReceiverRepository{Handler: handler}
}

func (receiver ReceiverRepository) Save(fileToConvert models.FileToConvert) (int, error) {

	// Append to the fileToConvert table
	if result := receiver.Handler.DB.Create(&fileToConvert); result.Error != nil {
		return 0, fmt.Errorf("error while inserting file: %v", result.Error)
	}
	return fileToConvert.Id, nil

}
