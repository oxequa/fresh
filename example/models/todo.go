package models

type Todo struct {
	Uuid string  `json:"uuid"`
	Title string `json:"title"`
	UserUuid string `json:"userUuid"`
}
