package model

type Email struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Header string `json:"header"`
	Body   string `json:"body"`
}
