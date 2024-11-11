package models

type Request struct {
	Title       string `json:"title"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}
