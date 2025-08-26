package services

import (
	"kafkapractisel0/repo"
)

type EmulatorService interface {
	GetRandomDelivery() (int, error)
}

type emulatorService struct {
	c repo.CustomerRepo
	d repo.DeliveryRepo
	i repo.ItemsRepo
}

func NewEmulatorService(c repo.CustomerRepo, d repo.DeliveryRepo, i repo.ItemsRepo) EmulatorService {
	return &emulatorService{c, d, i}
}

func (s *emulatorService) GetRandomDelivery() (int, error) {
	return s.d.GetRandomDelivery()
}
