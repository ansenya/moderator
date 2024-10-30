package models

type Request struct {
	Topic string `json:"topic"`
	Text  string `json:"text"`
}
