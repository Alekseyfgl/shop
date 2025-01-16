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

type nodeTypeRepository struct{}

type NodeTypeRepositoryInterface interface {
	GetAllNodeTypes(pageNumber, pageSize int) ([]model.NodeTypeRow, int, error)
	CreateNodeType(size *dto.CreateNodeTypeRequest) (int, error)
	UpdateNodeType(size *model.NodeTypeRow) error
	DeleteNodeTypeById(id int) error
	GetNodeTypeById(id int) (*model.NodeTypeRow, error)
}

func NewNodeTypeRepository() NodeTypeRepositoryInterface {
	return &nodeTypeRepository{}
}

var NodeTypeRepo = NewNodeTypeRepository()

func (r *nodeTypeRepository) GetAllNodeTypes(pageNumber, pageSize int) ([]model.NodeTypeRow, int, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := utils.CalculateOffset(pageNumber, pageSize)

	var totalCount int
	err := pg_conf.GetDB().QueryRow("SELECT COUNT(*) FROM shop.node_types").Scan(&totalCount)
	if err != nil {
		log.Error("Failed to count characteristic", zap.Error(err))
		return nil, 0, err
	}

	rows, err := pg_conf.GetDB().Query("SELECT id, type, description FROM shop.node_types ORDER BY id ASC LIMIT $1 OFFSET $2", pageSize, offset)
	if err != nil {
		log.Error("Failed to fetch node_types", zap.Error(err))
		return nil, 0, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Warn("Failed to close rows", zap.Error(closeErr))
		}
	}()

	scanFunc := func(rows *sql.Rows) (model.NodeTypeRow, error) {
		var size model.NodeTypeRow
		if err := rows.Scan(&size.ID, &size.Type, &size.Description); err != nil {
			return model.NodeTypeRow{}, err
		}
		return size, nil
	}

	sizes, err := utils.DecodeRows[model.NodeTypeRow](rows, scanFunc)
	if err != nil {
		log.Error("Failed to decode node_type", zap.Error(err))
		return nil, 0, err
	}

	return sizes, totalCount, nil
}

func (r *nodeTypeRepository) CreateNodeType(size *dto.CreateNodeTypeRequest) (int, error) {
	var insertedID int
	err := pg_conf.GetDB().QueryRow(
		"INSERT INTO shop.node_types (type, description) VALUES ($1, $2) RETURNING id",
		size.Type, size.Description,
	).Scan(&insertedID)

	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (r *nodeTypeRepository) UpdateNodeType(size *model.NodeTypeRow) error {
	_, err := pg_conf.GetDB().Exec(
		"UPDATE shop.node_types SET type = $1, description = $2 WHERE id = $3",
		size.Type, size.Description, size.ID,
	)
	return err
}

func (r *nodeTypeRepository) DeleteNodeTypeById(id int) error {
	_, err := pg_conf.GetDB().Exec("DELETE FROM shop.node_types WHERE id = $1", id)
	return err
}

func (r *nodeTypeRepository) GetNodeTypeById(id int) (*model.NodeTypeRow, error) {
	var size model.NodeTypeRow

	err := pg_conf.GetDB().QueryRow(
		"SELECT id, type, description FROM shop.node_types WHERE id = $1",
		id,
	).Scan(&size.ID, &size.Type, &size.Description)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("node_type not found", zap.Int("id", id))
			return nil, errors.New("node_type not found")
		}
		log.Error("Failed to fetch node_type by ID", zap.Error(err))
		return nil, err
	}

	return &size, nil
}
