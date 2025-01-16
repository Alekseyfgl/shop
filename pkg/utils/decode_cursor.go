package utils

import (
	"database/sql"
	"go.uber.org/zap"
	"shop/pkg/log"
)

func DecodeRows[T any](rows *sql.Rows, scanFunc func(*sql.Rows) (T, error)) ([]T, error) {
	var items []T
	for rows.Next() {
		item, err := scanFunc(rows)
		if err != nil {
			log.Error("Failed to scan row", zap.Error(err))
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		log.Error("Rows iteration encountered an error", zap.Error(err))
		return nil, err
	}

	return items, nil
}
