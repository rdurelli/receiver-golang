package services

import (
	"receiver/models"
	"receiver/repositories"
)

type ReceiverService struct {
	ReceiverRepository repositories.ReceiverRepository
}

func New(receiverRepository repositories.ReceiverRepository) ReceiverService {
	return ReceiverService{ReceiverRepository: receiverRepository}
}

func (receiverService ReceiverService) Save(file models.FileToConvert) (int, bool) {
	result, err := receiverService.ReceiverRepository.Save(file)

	if err != nil {
		return 0, true
	}
	return result, false
}
