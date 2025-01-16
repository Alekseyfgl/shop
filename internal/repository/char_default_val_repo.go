package repository

import (
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"shop/configs/pg_conf"
	"shop/internal/api/dto"
	"shop/internal/model"
	"shop/pkg/log"
	"shop/pkg/utils"
)

type charDefaultValueRepository struct{}

type CharDefaultValueInterface interface {
	GetAllDefaultValues(pageNumber, pageSize int) ([]model.CharDefaultValue, int, error)
	CreateDefaultValue(size *dto.CreateCharDefValueRequest) (int, error)
	UpdateDefaultValue(size *dto.UpdateCharDefValueRequest) error
	DeleteDefaultValueById(id int) error
	GetDefaultValueById(id int) (*model.CharDefaultValueRow, error)
}

func NewCharDefaultValueRepository() CharDefaultValueInterface {
	return &charDefaultValueRepository{}
}

var CharDefaultValueRepo = NewCharDefaultValueRepository()

func (r *charDefaultValueRepository) GetAllDefaultValues(pageNumber, pageSize int) ([]model.CharDefaultValue, int, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := utils.CalculateOffset(pageNumber, pageSize)

	var totalCount int
	err := pg_conf.GetDB().QueryRow("SELECT COUNT(*) FROM shop.char_default_value").Scan(&totalCount)
	if err != nil {
		log.Error("Failed to count char_default_value", zap.Error(err))
		return nil, 0, err
	}

	rows, err := pg_conf.GetDB().Query(`
			SELECT cdv.id,
				   cdv.characteristic_id,
				   cdv.value,
				   ch.title
			FROM shop.characteristics ch
					 JOIN shop.char_default_value cdv on ch.id = cdv.characteristic_id
			 ORDER BY cdv.id ASC LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		log.Error("Failed to fetch char_default_value", zap.Error(err))
		return nil, 0, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Warn("Failed to close rows", zap.Error(closeErr))
		}
	}()

	scanFunc := func(rows *sql.Rows) (model.CharDefaultValue, error) {
		var row model.CharDefaultValue
		if err := rows.Scan(
			&row.ID,
			&row.CharacteristicId,
			&row.Value,
			&row.Title,
		); err != nil {
			return model.CharDefaultValue{}, err
		}
		return row, nil
	}

	data, err := utils.DecodeRows[model.CharDefaultValue](rows, scanFunc)
	if err != nil {
		log.Error("Failed to decode char_default_value", zap.Error(err))
		return nil, 0, err
	}

	return data, totalCount, nil
}

func (r *charDefaultValueRepository) CreateDefaultValue(data *dto.CreateCharDefValueRequest) (int, error) {
	var insertedID int
	err := pg_conf.GetDB().QueryRow(
		"INSERT INTO shop.char_default_value (characteristic_id,value) VALUES ($1, $2) RETURNING id",
		data.CharacteristicId, data.Value,
	).Scan(&insertedID)

	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (r *charDefaultValueRepository) UpdateDefaultValue(data *dto.UpdateCharDefValueRequest) error {
	_, err := pg_conf.GetDB().Exec(
		"UPDATE shop.char_default_value SET value = $1 WHERE id = $2",
		data.Value, data.ID,
	)
	return err
}

func (r *charDefaultValueRepository) DeleteDefaultValueById(id int) error {
	_, err := pg_conf.GetDB().Exec(
		"DELETE FROM shop.char_default_value  WHERE id = $1",
		id,
	)
	if err != nil {
		// Log the error with context for easier debugging
		log.Error("Failed to delete node", zap.Int("id", id), zap.Error(err))
		return err
	}
	return nil
}

func (r *charDefaultValueRepository) GetDefaultValueById(id int) (*model.CharDefaultValueRow, error) {
	var data model.CharDefaultValueRow

	err := pg_conf.GetDB().QueryRow(
		"SELECT * FROM shop.char_default_value WHERE id = $1",
		id,
	).Scan(
		&data.ID,
		&data.CharacteristicId,
		&data.Value,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("char_default_value not found", zap.Int("id", id))
			return nil, errors.New("node not found")
		}
		log.Error("Failed to fetch char_default_value by ID", zap.Error(err))
		return nil, err
	}

	return &data, nil
}
