package store

import (
	"context"
	"go_carzone/models"
)

type CarStoreInterface interface {
	CreateCar(ctx context.Context,carReq *models.CarRequest) (models.Car,error)
	GetCarByID(ctx context.Context,id string) (models.Car,error)
	GetCarsByBrand(ctx context.Context,brand string,isEngine bool) ([]models.Car,error)
	UpdateCar(ctx context.Context,id string,carReq *models.CarRequest) (models.Car,error)
	DeleteCar(ctx context.Context,id string) (models.Car,error)
}