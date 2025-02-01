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
	"time"
)

type nodeRepository struct{}

type NodeRepositoryInterface interface {
	GetAllNodes(pageNumber, pageSize int) ([]model.NodeRow, int, error)
	CreateNode(size *dto.CreateNodeRequest) (int, error)
	UpdateNodes(size *dto.UpdateNodeRequest) error
	DeleteNodeById(id int) error
	GetNodeById(id int) (*model.NodeRow, error)
}

func NewNodeRepository() NodeRepositoryInterface {
	return &nodeRepository{}
}

var NodeRepo = NewNodeRepository()

func (r *nodeRepository) GetAllNodes(pageNumber, pageSize int) ([]model.NodeRow, int, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := utils.CalculateOffset(pageNumber, pageSize)

	var totalCount int
	err := pg_conf.GetDB().QueryRow("SELECT COUNT(*) FROM shop.nodes").Scan(&totalCount)
	if err != nil {
		log.Error("Failed to count characteristic", zap.Error(err))
		return nil, 0, err
	}

	rows, err := pg_conf.GetDB().Query("SELECT id, title, node_type_id, description, created_at, updated_at, removed_at FROM shop.nodes ORDER BY id ASC LIMIT $1 OFFSET $2", pageSize, offset)
	if err != nil {
		log.Error("Failed to fetch node_types", zap.Error(err))
		return nil, 0, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Warn("Failed to close rows", zap.Error(closeErr))
		}
	}()

	scanFunc := func(rows *sql.Rows) (model.NodeRow, error) {
		var node model.NodeRow
		if err := rows.Scan(
			&node.ID,
			&node.Title,
			&node.NodeTypeId,
			&node.Description,
			&node.CreatedAt,
			&node.UpdatedAt,
			&node.RemovedAt,
		); err != nil {
			return model.NodeRow{}, err
		}
		return node, nil
	}

	nodes, err := utils.DecodeRows[model.NodeRow](rows, scanFunc)
	if err != nil {
		log.Error("Failed to decode node", zap.Error(err))
		return nil, 0, err
	}

	return nodes, totalCount, nil
}

func (r *nodeRepository) CreateNode(node *dto.CreateNodeRequest) (int, error) {
	var insertedID int
	err := pg_conf.GetDB().QueryRow(
		"INSERT INTO shop.nodes (title,node_type_id, description) VALUES ($1, $2, $3) RETURNING id",
		node.Title, node.NodeTypeId, node.Description,
	).Scan(&insertedID)

	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (r *nodeRepository) UpdateNodes(node *dto.UpdateNodeRequest) error {
	_, err := pg_conf.GetDB().Exec(
		"UPDATE shop.nodes SET title = $1, node_type_id = $2, description = $3 WHERE id = $4",
		node.Title, node.NodeTypeId, node.Description, node.ID,
	)
	return err
}

func (r *nodeRepository) DeleteNodeById(id int) error {
	// Get the current time in UTC
	currentTime := time.Now().UTC()

	// Execute the UPDATE statement with currentTime and id as parameters
	_, err := pg_conf.GetDB().Exec(
		"UPDATE shop.nodes SET removed_at = $1 WHERE id = $2",
		currentTime,
		id,
	)
	if err != nil {
		// Log the error with context for easier debugging
		log.Error("Failed to delete node", zap.Int("id", id), zap.Error(err))
		return err
	}
	return nil
}
func (r *nodeRepository) GetNodeById(id int) (*model.NodeRow, error) {
	var node model.NodeRow

	err := pg_conf.GetDB().QueryRow(
		"SELECT id, title, node_type_id, description, created_at, updated_at, removed_at FROM shop.nodes WHERE id = $1",
		id,
	).Scan(
		&node.ID,
		&node.Title,
		&node.NodeTypeId,
		&node.Description,
		&node.CreatedAt,
		&node.UpdatedAt,
		&node.RemovedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("Node not found", zap.Int("id", id))
			return nil, errors.New("node not found")
		}
		log.Error("Failed to fetch node by ID", zap.Error(err))
		return nil, err
	}

	return &node, nil
}
