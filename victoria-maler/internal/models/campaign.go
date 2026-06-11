package models

type Campaign struct {
    ID          int64  `db:"id" json:"id"`
    Title       string `db:"title" json:"title"`
    Description string `db:"description" json:"description"`
}
