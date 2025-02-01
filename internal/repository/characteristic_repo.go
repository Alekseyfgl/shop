package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"shop/configs/pg_conf"
	"shop/internal/api/dto"
	"shop/internal/model"
	"shop/pkg/log"
	"shop/pkg/utils"
)

type characteristicRepository struct{}

type CharacteristicRepositoryInterface interface {
	GetAllCharacteristics(pageNumber, pageSize int) ([]model.CharacteristicRow, int, error)
	CreateCharacteristics(size *dto.CreateCharacteristicRequest) (int, error)
	UpdateCharacteristics(size *model.CharacteristicRow) error
	DeleteCharacteristicsById(id int) error
	GetCharacteristicsById(id int) (*model.CharacteristicRow, error)
	CheckCharsByIds(ids []int) error
	GetCharFilters() (*[]model.CharFiltersRow, error)
}

func NewCharacteristicRepository() CharacteristicRepositoryInterface {
	return &characteristicRepository{}
}

var CharacteristicRepo = NewCharacteristicRepository()

func (r *characteristicRepository) GetAllCharacteristics(pageNumber, pageSize int) ([]model.CharacteristicRow, int, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := utils.CalculateOffset(pageNumber, pageSize)

	var totalCount int
	err := pg_conf.GetDB().QueryRow("SELECT COUNT(*) FROM shop.characteristics").Scan(&totalCount)
	if err != nil {
		log.Error("Failed to count characteristic", zap.Error(err))
		return nil, 0, err
	}

	rows, err := pg_conf.GetDB().Query("SELECT id, title, description FROM shop.characteristics ORDER BY id ASC LIMIT $1 OFFSET $2", pageSize, offset)
	if err != nil {
		log.Error("Failed to fetch characteristics", zap.Error(err))
		return nil, 0, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Warn("Failed to close rows", zap.Error(closeErr))
		}
	}()

	scanFunc := func(rows *sql.Rows) (model.CharacteristicRow, error) {
		var char model.CharacteristicRow
		if err := rows.Scan(&char.ID, &char.Title, &char.Description); err != nil {
			return model.CharacteristicRow{}, err
		}
		return char, nil
	}

	chars, err := utils.DecodeRows[model.CharacteristicRow](rows, scanFunc)
	if err != nil {
		log.Error("Failed to decode characteristics", zap.Error(err))
		return nil, 0, err
	}

	if chars == nil {
		chars = make([]model.CharacteristicRow, 0)
	}

	return chars, totalCount, nil
}

func (r *characteristicRepository) CreateCharacteristics(data *dto.CreateCharacteristicRequest) (int, error) {
	var insertedID int
	err := pg_conf.GetDB().QueryRow(
		"INSERT INTO shop.characteristics (title, description) VALUES ($1, $2) RETURNING id",
		data.Title, data.Description,
	).Scan(&insertedID)

	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (r *characteristicRepository) UpdateCharacteristics(data *model.CharacteristicRow) error {
	_, err := pg_conf.GetDB().Exec(
		"UPDATE shop.characteristics SET title = $1, description = $2 WHERE id = $3",
		data.Title, data.Description, data.ID,
	)
	return err
}

func (r *characteristicRepository) DeleteCharacteristicsById(id int) error {
	_, err := pg_conf.GetDB().Exec("DELETE FROM shop.characteristics WHERE id = $1", id)
	return err
}

func (r *characteristicRepository) GetCharacteristicsById(id int) (*model.CharacteristicRow, error) {
	var char model.CharacteristicRow

	err := pg_conf.GetDB().QueryRow(
		"SELECT id, title, description FROM shop.characteristics WHERE id = $1",
		id,
	).Scan(&char.ID, &char.Title, &char.Description)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("Title not found", zap.Int("id", id))
			return nil, errors.New("size not found")
		}
		log.Error("Failed to fetch characteristic by ID", zap.Error(err))
		return nil, err
	}

	return &char, nil
}

func (r *characteristicRepository) CheckCharsByIds(ids []int) error {
	// Выполняем запрос, чтобы получить все характеристики с указанными id
	rows, err := pg_conf.GetDB().Query(
		"SELECT id FROM shop.characteristics WHERE id = ANY($1)",
		pq.Array(ids),
	)
	if err != nil {
		log.Error("Failed to fetch characteristics by IDs", zap.Error(err))
		return err
	}
	defer rows.Close()

	// Сохраняем найденные ID
	foundIDs := make(map[int]struct{})
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Error("Failed to scan characteristic row", zap.Error(err))
			return err
		}
		foundIDs[id] = struct{}{}
	}

	// Проверяем наличие ошибок при обработке строк
	if err := rows.Err(); err != nil {
		log.Error("Error occurred during rows iteration", zap.Error(err))
		return err
	}

	// Определяем отсутствующие ID
	missingIDs := []int{}
	for _, id := range ids {
		if _, exists := foundIDs[id]; !exists {
			missingIDs = append(missingIDs, id)
		}
	}

	// Если есть отсутствующие ID, возвращаем ошибку
	if len(missingIDs) > 0 {
		log.Warn("Some IDs are missing", zap.Ints("missingIDs", missingIDs))
		return fmt.Errorf("these IDs don't  exist: %v", missingIDs)
	}

	return nil
}

func (r *characteristicRepository) GetCharFilters() (*[]model.CharFiltersRow, error) {
	// Выполняем запрос, чтобы получить все характеристики с указанными id
	rows, err := pg_conf.GetDB().Query(
		`
			SELECT ch.id  "characteristicId",
				   ch.title,
				   ch.description,
				   cdv.value
			FROM shop.characteristics  ch
			LEFT JOIN shop.char_default_value cdv on ch.id = cdv.characteristic_id
			WHERE ch.is_visible = true`,
	)

	if err != nil {
		log.Error("Failed to fetch filters", zap.Error(err))
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Warn("Failed to close rows", zap.Error(closeErr))
		}
	}()

	scanFunc := func(rows *sql.Rows) (model.CharFiltersRow, error) {
		var filter model.CharFiltersRow
		if err := rows.Scan(&filter.CharacteristicId, &filter.Title, &filter.Description, &filter.Value); err != nil {
			return model.CharFiltersRow{}, err
		}
		return filter, nil
	}

	filters, err := utils.DecodeRows[model.CharFiltersRow](rows, scanFunc)
	if err != nil {
		log.Error("Failed to decode filters", zap.Error(err))
		return nil, err
	}

	return &filters, nil
}
