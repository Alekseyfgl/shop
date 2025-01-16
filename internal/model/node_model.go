package model

type NodeRow struct {
	ID          int     `db:"id" json:"id"`
	Title       string  `db:"title" json:"title"`
	NodeTypeId  int     `db:"node_type_id" json:"nodeTypeId"`
	Description *string `db:"description" json:"description"`
	CreatedAt   string  `db:"created_at" json:"createdAt"`
	UpdatedAt   string  `db:"updated_at" json:"updatedAt"`
	RemovedAt   *string `db:"removed_at" json:"removedAt"`
}
