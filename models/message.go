package models

type Message struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Query    string `json:"query"`
	Response string `json:"response"`
}
