package model

type Email struct {
	FromName   string `json:"from_name"`
	ToEmail     string `json:"to_email"`
	Header string `json:"header"`
	Body   string `json:"body"`
}
