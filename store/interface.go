package store

import (
	"context"
	"github.com/iangechuki/go_carzone/models"
)

type CarStoreInterface interface {
	CreateCar(ctx context.Context,carReq *models.CarRequest) (models.Car,error)
	GetCarByID(ctx context.Context,id string) (models.Car,error)
	GetCarsByBrand(ctx context.Context,brand string,isEngine bool) ([]models.Car,error)
	UpdateCar(ctx context.Context,id string,carReq *models.CarRequest) (models.Car,error)
	DeleteCar(ctx context.Context,id string) (models.Car,error)
}

type EngineStoreInterface interface {
	GetEngineByID(ctx context.Context,id string) (models.Engine,error)
	CreateEngine(ctx context.Context,engineReq *models.EngineRequest) (models.Engine,error)
	UpdateEngine(ctx context.Context,id string,engineReq *models.EngineRequest) (models.Engine,error)
	DeleteEngine(ctx context.Context,id string) (models.Engine,error)
}