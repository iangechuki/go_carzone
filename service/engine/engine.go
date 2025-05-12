package engine

import (
	"context"
	"go_carzone/models"
	"go_carzone/store"
)

type EngineService struct{
	store store.EngineStoreInterface
}

func NewEngineService(store store.EngineStoreInterface) *EngineService {
	return &EngineService{
		store: store,
	}
}

func (s *EngineService)GetEngineByID(ctx context.Context,id string) (*models.Engine,error) {
	engine ,err := s.store.GetEngineByID(ctx,id)
	if err != nil {
		return nil,err
	}
	return &engine,nil
}

func (s *EngineService)CreateEngine(ctx context.Context,engineReq *models.EngineRequest) (*models.Engine,error) {
	if err := models.ValidateEngineRequest(*engineReq); err != nil {
		return nil,err
	}
	engine ,err := s.store.CreateEngine(ctx,engineReq)
	if err != nil {
		return nil,err
	}
	return &engine,nil
}

func (s *EngineService)UpdateEngine(ctx context.Context,id string,engineReq *models.EngineRequest) (*models.Engine,error) {
	if err := models.ValidateEngineRequest(*engineReq); err != nil {
		return nil,err
	}
	engine ,err := s.store.UpdateEngine(ctx,id,engineReq)
	if err != nil {
		return nil,err
	}
	return &engine,nil
}

func (s *EngineService)DeleteEngine(ctx context.Context,id string) (*models.Engine,error) {
	engine ,err := s.store.DeleteEngine(ctx,id)
	if err != nil {
		return nil,err
	}
	return &engine,nil
}