package model

type CharDefaultValueRow struct {
	ID               int    `db:"id" json:"id"`
	CharacteristicId int    `db:"characteristicId" json:"characteristicId"`
	Value            string `db:"value" json:"value"`
}

type CharDefaultValue struct {
	CharDefaultValueRow
	Title string `db:"title" json:"title"`
}
