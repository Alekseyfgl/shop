package model

type CharacteristicRow struct {
	ID          int     `db:"id" json:"id"`
	Title       string  `db:"title" json:"title"`
	Description *string `db:"description" json:"description"`
	IsVisible   bool    `db:"is_visible" json:"isVisible"`
}

type CharFiltersRow struct {
	CharacteristicId int     `db:"characteristicId" json:"characteristicId"`
	Title            string  `db:"title" json:"title"`
	Description      *string `db:"description" json:"description"`
	Value            *string `db:"value" json:"value"`
}

type CharFilterResponse struct {
	CharacteristicId int      `db:"characteristicId" json:"characteristicId"`
	Title            string   `db:"title" json:"title"`
	Values           []string `db:"values" json:"values"`
}
