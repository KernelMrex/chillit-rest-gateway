package models

type Place struct {
	Id          uint64 `json:"id"`
	Title       string `json:"title"`
	Address     string `json:"address"`
	Description string `json:"description"`
}
