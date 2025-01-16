package model

type SizeRow struct {
	ID          int     `db:"id" json:"id"`
	Title       string  `db:"title" json:"title"`
	Description *string `db:"description" json:"description"`
}
