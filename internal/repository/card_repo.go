package repository

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"shop/configs/pg_conf"
	"shop/internal/api/dto"
	"shop/internal/model"
	"shop/pkg/log"
	"shop/pkg/utils"
	"strconv"
	"strings"
)

type cardRepository struct{}

// CardRepositoryInterface описывает методы, необходимые для работы с "карточками" (nodes).
type CardRepositoryInterface interface {
	GetCardById(id int) (*[]model.CardRow, error)
	GetAllCards(pageNumber, pageSize int, filters *[]model.CardFilter) (*[]model.CardRow, int, error)
	CreateCard(dto *dto.CreateCardDTO) (int, error)
	FindByVectorSearch(text string, limit int) (*[]model.CardRow, error)
}

// NewCardRepository создаёт новый экземпляр репозитория для работы с карточками.
func NewCardRepository() CardRepositoryInterface {
	return &cardRepository{}
}

var CardRepo = NewCardRepository()

func (r *cardRepository) FindByVectorSearch(text string, limit int) (*[]model.CardRow, error) {
	if limit < 1 || limit > 20 {
		limit = 10
	}

	query := `
        SELECT
            n.id           AS "nodeId",
            n.title,
            n.description  AS "nodeDescription",
            n.created_at   AS "createdAt",
            n.updated_at   AS "updatedAt",
            n.removed_at   AS "removedAt",
            COALESCE(string_to_array(n.images, ','), '{}') AS "images",
            nt.type        AS "nodeType",
            nt.description AS "nodeTypeDescription",
            c.title        AS "characteristic",
            cv.value       AS "characteristicValue",
            cv.add_params  AS "additionalParams",
            c.description  AS "characteristicDescription"
        FROM shop.nodes n
                 JOIN shop.node_types nt ON nt.id = n.node_type_id
                 JOIN shop.characteristic_values cv ON n.id = cv.node_id
                 JOIN shop.characteristics c ON c.id = cv.characteristic_id
        WHERE n.search_vector @@ plainto_tsquery('russian', $1)
        ORDER BY ts_rank_cd(n.search_vector, plainto_tsquery('russian', $1)) DESC
        LIMIT $2
    `

	rows, err := pg_conf.GetDB().Query(query, text, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []model.CardRow
	for rows.Next() {
		var card model.CardRow
		if err := rows.Scan(
			&card.NodeId,
			&card.Title,
			&card.NodeDescription,
			&card.CreatedAt,
			&card.UpdatedAt,
			&card.RemovedAt,
			pq.Array(&card.Images),
			&card.NodeType,
			&card.NodeTypeDescription,
			&card.Characteristic,
			&card.CharacteristicValue,
			&card.AdditionalParams,
			&card.CharacteristicDescription,
		); err != nil {
			log.Error("Failed to scan row in FindByVectorSearch", zap.Error(err))
			return nil, err
		}
		cards = append(cards, card)
	}

	if err := rows.Err(); err != nil {
		log.Error("Row iteration error in FindByVectorSearch", zap.Error(err))
		return nil, err
	}

	return &cards, nil
}

// GetCardById возвращает список характеристик (CardRow) для заданного nodeId.
func (r *cardRepository) GetCardById(id int) (*[]model.CardRow, error) {
	// Выполняем запрос к базе данных
	rows, err := pg_conf.GetDB().Query(
		`
        SELECT n.id           AS "nodeId",
               n.title,
               n.description  AS "nodeDescription",
               n.created_at   AS "createdAt",
               n.updated_at   AS "updatedAt",
               n.removed_at   AS "removedAt",
               COALESCE(string_to_array(n.images, ','), '{}') AS "images",
               nt.type        AS "nodeType",
               nt.description AS "nodeTypeDescription",
               c.title        AS characteristic,
               cv.value       AS "characteristicValue",
               cv.add_params  AS "additionalParams",
               c.description  AS "characteristicDescription"
        FROM shop.nodes n
                 JOIN shop.node_types nt ON nt.id = n.node_type_id
                 JOIN shop.characteristic_values cv ON n.id = cv.node_id
                 JOIN shop.characteristics c ON c.id = cv.characteristic_id
        WHERE n.id = $1
        `,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []model.CardRow

	// Обрабатываем строки результата
	for rows.Next() {
		var card model.CardRow
		if err := rows.Scan(
			&card.NodeId,
			&card.Title,
			&card.NodeDescription,
			&card.CreatedAt,
			&card.UpdatedAt,
			&card.RemovedAt,
			pq.Array(&card.Images),
			&card.NodeType,
			&card.NodeTypeDescription,
			&card.Characteristic,
			&card.CharacteristicValue,
			&card.AdditionalParams,
			&card.CharacteristicDescription,
		); err != nil {
			log.Error("Failed to scan row", zap.Error(err))
			return nil, err
		}
		cards = append(cards, card)
	}

	// Проверяем ошибки, возникшие при итерации
	if err := rows.Err(); err != nil {
		log.Error("Error during rows iteration", zap.Error(err))
		return nil, err
	}

	return &cards, nil
}

func (r *cardRepository) GetAllCards(pageNumber, pageSize int, filters *[]model.CardFilter) (*[]model.CardRow, int, error) {
	// Установка значений по умолчанию для пагинации
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := utils.CalculateOffset(pageNumber, pageSize)

	// Генерируем часть WHERE на основе фильтров и получаем аргументы.
	whereClause, whereArgs := buildWhereClause(*filters)

	// -----------------------------------------------------------
	// Считаем количество с учётом фильтров
	// -----------------------------------------------------------
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM (
			SELECT n.id
			FROM shop.nodes n
			         JOIN shop.node_types nt ON nt.id = n.node_type_id
			         JOIN shop.characteristic_values cv ON n.id = cv.node_id
			         JOIN shop.characteristics c ON c.id = cv.characteristic_id
			%s
			GROUP BY n.id
		) grouped_nodes;
	`, whereClause)

	var totalCount int
	if err := pg_conf.GetDB().QueryRow(countQuery, whereArgs...).Scan(&totalCount); err != nil {
		log.Error("Failed to count cards with filters", zap.Error(err))
		return nil, 0, err
	}

	// -----------------------------------------------------------
	// Формируем основной SELECT с пагинацией
	// -----------------------------------------------------------
	args := make([]interface{}, 0, len(whereArgs)+2)
	args = append(args, whereArgs...)
	args = append(args, pageSize)
	args = append(args, offset)

	selectQuery := fmt.Sprintf(`
        SELECT n.id           AS "nodeId",
               n.title,
               n.description  AS "nodeDescription",
               n.created_at   AS "createdAt",
               n.updated_at   AS "updatedAt",
               n.removed_at   AS "removedAt",
               COALESCE(string_to_array(n.images, ','), '{}') AS "images",
               nt.id          AS "nodeType", -- изменено с nt.type на nt.id
               nt.description AS "nodeTypeDescription",
               c.title        AS "characteristic",
               cv.value       AS "characteristicValue",
               cv.add_params  AS "additionalParams",
               c.description  AS "characteristicDescription"
        FROM shop.nodes n
                 JOIN shop.node_types nt ON nt.id = n.node_type_id
                 JOIN shop.characteristic_values cv ON n.id = cv.node_id
                 JOIN shop.characteristics c ON c.id = cv.characteristic_id
        %s
        ORDER BY n.created_at DESC
        LIMIT $%d OFFSET $%d
    `,
		whereClause,      // подставляем строку WHERE (может быть пустой, если фильтров нет)
		len(whereArgs)+1, // placeholder для LIMIT
		len(whereArgs)+2, // placeholder для OFFSET
	)

	rows, err := pg_conf.GetDB().Query(selectQuery, args...)
	if err != nil {
		log.Error("Failed to fetch cards", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	cards := make([]model.CardRow, 0, pageSize)
	for rows.Next() {
		var card model.CardRow
		if err := rows.Scan(
			&card.NodeId,
			&card.Title,
			&card.NodeDescription,
			&card.CreatedAt,
			&card.UpdatedAt,
			&card.RemovedAt,
			pq.Array(&card.Images),
			&card.NodeType,
			&card.NodeTypeDescription,
			&card.Characteristic,
			&card.CharacteristicValue,
			&card.AdditionalParams,
			&card.CharacteristicDescription,
		); err != nil {
			log.Error("Failed to scan row", zap.Error(err))
			return nil, 0, err
		}
		cards = append(cards, card)
	}

	if err := rows.Err(); err != nil {
		log.Error("Error during rows iteration", zap.Error(err))
		return nil, 0, err
	}

	return &cards, totalCount, nil
}

// buildWhereClause динамически формирует часть WHERE с placeholder’ами.
// Пример фильтров:
//
//	[
//	    { Key: "Размеры",   Values: "M" },
//	    { Key: "Скидка",    Values: "" },
//	    { Key: "nodeType",  Values: 1 } // значение виде int
//	]
//
// Если присутствует фильтр nodeType, то условие будет вида (nt.id = $X)
func buildWhereClause(filters []model.CardFilter) (string, []interface{}) {
	if len(filters) == 0 {
		return "", nil
	}

	var (
		conditions       []string
		args             []interface{}
		placeholderIndex = 1
	)

	for _, f := range filters {
		// Обработка специального фильтра nodeType
		if f.Key == "nodeTypeId" {
			// Приводим значение к int. Если конвертация не удалась, фильтр пропускается.
			nodeTypeID, err := strconv.Atoi(f.Values)
			if err != nil {
				log.Warn("Invalid nodeType filter value, skipping nodeType filter", zap.String("value", f.Values))
				continue
			}
			conditions = append(conditions, fmt.Sprintf("(nt.id = $%d)", placeholderIndex))
			args = append(args, nodeTypeID)
			placeholderIndex++
			continue
		}

		// Остальные фильтры
		if f.Values == "" {
			// Условие вида (c.title = $X)
			conditions = append(conditions, fmt.Sprintf("(c.title = $%d)", placeholderIndex))
			args = append(args, f.Key)
			placeholderIndex++
		} else {
			// Условие вида (c.title = $X AND cv.value = $%d)
			conditions = append(conditions, fmt.Sprintf("(c.title = $%d AND cv.value = $%d)", placeholderIndex, placeholderIndex+1))
			args = append(args, f.Key, f.Values)
			placeholderIndex += 2
		}
	}

	if len(conditions) == 0 {
		return "", nil
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")
	return whereClause, args
}

// CreateCard реализует логику создания node и его характеристик.
func (r *cardRepository) CreateCard(dto *dto.CreateCardDTO) (int, error) {
	// Начинаем транзакцию
	tx, err := pg_conf.GetDB().Begin()
	if err != nil {
		log.Error("Failed to begin transaction", zap.Error(err))
		return 0, err
	}

	defer func() {
		// Если случится паника — откатываемся
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// 1. Вставляем запись в shop.nodes
	newNodeID, err := r.insertNodeTx(tx, dto)
	if err != nil {
		_ = tx.Rollback()
		log.Error("Failed to insert node", zap.Error(err))
		return 0, err
	}

	// 2. Вставляем все характеристики (bulk insert)
	err = r.insertCharacteristicsTx(tx, newNodeID, dto.Characteristics)
	if err != nil {
		_ = tx.Rollback()
		log.Error("Failed to insert characteristics", zap.Error(err))
		return 0, err
	}

	// 3. Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		log.Error("Failed to commit transaction", zap.Error(err))
		return 0, err
	}

	log.Info(fmt.Sprintf("Card successfully created with ID: %d", newNodeID))
	return newNodeID, nil
}

// insertNodeTx — вспомогательная функция, вставляет запись в shop.nodes.
func (r *cardRepository) insertNodeTx(tx *sql.Tx, dto *dto.CreateCardDTO) (int, error) {
	// Преобразуем массив изображений в строку через запятую
	imagesString := strings.Join(dto.Images, ",")
	const query = `
		INSERT INTO shop.nodes (title, description, node_type_id, images)
		VALUES ($1, $2, $3,$4)
		RETURNING id
	`

	var newNodeID int
	err := tx.QueryRow(query, dto.Title, dto.NodeDescription, dto.NodeTypeId, imagesString).Scan(&newNodeID)
	if err != nil {
		return 0, err
	}
	return newNodeID, nil
}

// insertCharacteristicsTx — вспомогательная функция, вставляет записи в shop.characteristic_values (bulk insert).
func (r *cardRepository) insertCharacteristicsTx(
	tx *sql.Tx,
	nodeID int,
	characteristics []dto.CharDTO,
) error {
	// Если нет характеристик — ничего не вставляем
	if len(characteristics) == 0 {
		return nil
	}

	const baseQuery = `
		INSERT INTO shop.characteristic_values (node_id, characteristic_id, add_params, value)
		VALUES
	`

	// Подготавливаем слайсы для плейсхолдеров и аргументов
	valuesPlaceholder := make([]string, 0, len(characteristics))
	args := make([]interface{}, 0, len(characteristics)*4)

	// Формируем bulk insert
	for i, ch := range characteristics {
		// Если AdditionalParams != nil, используем как есть (сырые JSON-байты), иначе nil
		var addParamsJSON interface{}
		if ch.AdditionalParams != nil {
			addParamsJSON = ch.AdditionalParams
		} else {
			addParamsJSON = nil
		}

		// Формируем плейсхолдеры для текущей итерации
		// Например, ($1, $2, $3, $4), затем ($5, $6, $7, $8), и так далее
		startIdx := i*4 + 1
		placeholder := fmt.Sprintf("($%d, $%d, $%d, $%d)",
			startIdx,   // node_id
			startIdx+1, // characteristic_id
			startIdx+2, // add_params
			startIdx+3, // value
		)
		valuesPlaceholder = append(valuesPlaceholder, placeholder)

		// Дополняем args
		args = append(args, nodeID, ch.Id, addParamsJSON, ch.Value)
	}

	// Склеиваем плейсхолдеры в один INSERT
	query := baseQuery + strings.Join(valuesPlaceholder, ",")
	_, err := tx.Exec(query, args...)
	return err
}
