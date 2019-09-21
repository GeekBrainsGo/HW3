package main

type BlogItems []BlogItem

type BlogItem struct {
	ID       int64
	Title    string
	Body     string
	Comments []string
}

// ErrorModel - Ошибка отвечаемая сервером
type ErrorModel struct {
	Code     int         `json:"code"`
	Err      string      `json:"error"`
	Desc     string      `json:"desc"`
	Internal interface{} `json:"internal"`
}
