package model

type NodeTypeRow struct {
	ID          int     `db:"id" json:"id"`
	Type        string  `db:"type" json:"type"`
	Description *string `db:"description" json:"description"`
}
