package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"shop/internal/api/dto"

	"go.uber.org/zap"
	"shop/configs/pg_conf"
	"shop/internal/model"
	"shop/pkg/log"
	"shop/pkg/utils"
)

type sizeRepository struct{}

type SizeRepositoryInterface interface {
	GetAllSizes(pageNumber, pageSize int) ([]model.SizeRow, int, error)
	CreateSize(size *dto.CreateSizeRequest) (int, error)
	UpdateSize(size *model.SizeRow) error
	DeleteSizeById(id int) error
	GetSizeById(id int) (*model.SizeRow, error)
}

func newSizeRepository() SizeRepositoryInterface {
	return &sizeRepository{}
}

var SizeRepo = newSizeRepository()

func (r *sizeRepository) GetAllSizes(pageNumber, pageSize int) ([]model.SizeRow, int, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := utils.CalculateOffset(pageNumber, pageSize)

	var totalCount int
	err := pg_conf.GetDB().QueryRow("SELECT COUNT(*) FROM size").Scan(&totalCount)
	if err != nil {
		log.Error("Failed to count sizes", zap.Error(err))
		return nil, 0, err
	}

	rows, err := pg_conf.GetDB().Query("SELECT id, title, description FROM size ORDER BY title DESC LIMIT $1 OFFSET $2", pageSize, offset)
	if err != nil {
		log.Error("Failed to fetch sizes", zap.Error(err))
		return nil, 0, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Warn("Failed to close rows", zap.Error(closeErr))
		}
	}()

	scanFunc := func(rows *sql.Rows) (model.SizeRow, error) {
		var size model.SizeRow
		if err := rows.Scan(&size.ID, &size.Title, &size.Description); err != nil {
			return model.SizeRow{}, err
		}
		return size, nil
	}
	fmt.Printf("scanFunc: %+v\n", scanFunc)

	sizes, err := utils.DecodeRows[model.SizeRow](rows, scanFunc)
	if err != nil {
		log.Error("Failed to decode sizes", zap.Error(err))
		return nil, 0, err
	}

	return sizes, totalCount, nil
}

func (r *sizeRepository) CreateSize(size *dto.CreateSizeRequest) (int, error) {
	var insertedID int
	err := pg_conf.GetDB().QueryRow(
		"INSERT INTO size (title, description) VALUES ($1, $2) RETURNING id",
		size.Title, size.Description,
	).Scan(&insertedID)

	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (r *sizeRepository) UpdateSize(size *model.SizeRow) error {
	_, err := pg_conf.GetDB().Exec(
		"UPDATE size SET size = $1, description = $2 WHERE id = $3",
		size.Title, size.Description, size.ID,
	)
	return err
}

func (r *sizeRepository) DeleteSizeById(id int) error {
	_, err := pg_conf.GetDB().Exec("DELETE FROM size WHERE id = $1", id)
	return err
}

func (r *sizeRepository) GetSizeById(id int) (*model.SizeRow, error) {
	var size model.SizeRow

	err := pg_conf.GetDB().QueryRow(
		"SELECT id, title, description FROM size WHERE id = $1",
		id,
	).Scan(&size.ID, &size.Title, &size.Description)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("Title not found", zap.Int("id", id))
			return nil, errors.New("size not found")
		}
		log.Error("Failed to fetch size by ID", zap.Error(err))
		return nil, err
	}

	return &size, nil
}
