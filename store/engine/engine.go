package engine

import (
	"context"
	"database/sql"
	"go_carzone/models"

	"github.com/google/uuid"
)

type EngineStore struct {
	db *sql.DB
}

func New(db *sql.DB) EngineStore {
	return EngineStore{
		db: db,
	}
}

func (s *EngineStore) GetEngineByID(ctx context.Context,id uuid.UUID) (models.Engine,error) {
	return models.Engine{},nil
}

func (s *EngineStore) CreateEngine(ctx context.Context,engineReq models.EngineRequest) (models.Engine,error) {
	return models.Engine{},nil
}

func (s *EngineStore) UpdateEngine(ctx context.Context,id uuid.UUID,engineReq *models.EngineRequest) (models.Engine,error) {
	return models.Engine{},nil
}

func (s *EngineStore) DeleteEngine(ctx context.Context,id uuid.UUID) (models.Engine,error) {
	return models.Engine{},nil
}