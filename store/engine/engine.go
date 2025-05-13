package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iangechuki/go_carzone/models"

	"github.com/google/uuid"
)

type EngineStore struct {
	db *sql.DB
}

func New(db *sql.DB) *EngineStore {
	return &EngineStore{
		db: db,
	}
}

func (s *EngineStore) GetEngineByID(ctx context.Context,id string) (models.Engine,error) {
	var engine models.Engine
	tx,err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{},err
	}
	defer func(){
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("could not rollback transaction: %v\n", rbErr)
			}
		}else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("could not commit transaction: %v\n", cmErr)
			}
		}
	}()
	err = tx.QueryRowContext(ctx,"SELECT id,displacement,no_of_cylinders,car_range FROM engine WHERE id = $1",id).
	Scan(
		&engine.EngineID,
		&engine.Displacement,
		&engine.NoOfCylinders,
		&engine.CarRange,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.Engine{},errors.New("engine not found")
		default:
			return models.Engine{},err
		}
	}
	return engine,nil
}

func (s *EngineStore) CreateEngine(ctx context.Context,engineReq *models.EngineRequest) (models.Engine,error) {
	tx,err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{},err
	}
	defer func ()  {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("could not rollback transaction: %v\n", rbErr)
			}
		}else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("could not commit transaction: %v\n", cmErr)
			}
		}
	}()
	engineID := uuid.New()

	_,err = tx.ExecContext(ctx,"INSERT INTO engine (id,displacement,no_of_cylinders,car_range,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)",
		engineID,
		engineReq.Displacement,
		engineReq.NoOfCylinders,
		engineReq.CarRange,
	)
	if err != nil {
		return models.Engine{},err
	}
	engine := models.Engine{
		EngineID: engineID,
		Displacement: engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange: engineReq.CarRange,

	}
	return engine,nil
}

func (s *EngineStore) UpdateEngine(ctx context.Context,id string,engineReq *models.EngineRequest) (models.Engine,error) {
	engineID,err := uuid.Parse(id)
	if err != nil {
		return models.Engine{},err
	}
	tx ,err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{},err
	}
	defer func () {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("could not rollback transaction: %v\n", rbErr)
			}
		}else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("could not commit transaction: %v\n", cmErr)
			}
		}
	}()
	results,err := tx.ExecContext(ctx,"UPDATE engine SET displacement = $1,no_of_cylinders = $2,car_range = $3 WHERE id = $4",
		engineReq.Displacement,
		engineReq.NoOfCylinders,
		engineReq.CarRange,
	)
	if err != nil {
		return models.Engine{},err
	}
	rowsAffected,err := results.RowsAffected()
	if err != nil {
		return models.Engine{},err
	}
	if rowsAffected == 0 {
		return models.Engine{},errors.New("no rows updated")
	}
	engine := models.Engine{
		EngineID: engineID,
		Displacement: engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange: engineReq.CarRange,
	}
	return engine,nil
}

func (s *EngineStore) DeleteEngine(ctx context.Context,id string) (models.Engine,error) {
	var engine models.Engine
	tx,err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{},err
	}
	defer func(){
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("could not rollback transaction: %v\n", rbErr)
			}
		}else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("could not commit transaction: %v\n", cmErr)
			}
		}
	}()

	err = tx.QueryRowContext(ctx,"SELECT id,displacement,no_of_cylinders,car_range FROM engine WHERE id = $1",id).Scan(
		&engine.EngineID,
		&engine.Displacement,
		&engine.NoOfCylinders,
		&engine.CarRange,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.Engine{},errors.New("engine not found")
		default:
			return models.Engine{},err
		}
	}
	result,err := tx.ExecContext(ctx,"DELETE FROM engine WHERE id = $1",id)
	if err != nil {
		return models.Engine{},err
	}
	rowsAffected,err := result.RowsAffected()
	if err != nil {
		return models.Engine{},err
	}
	if rowsAffected == 0 {
		return models.Engine{},errors.New("no rows deleted")
	}
	return engine,nil
}