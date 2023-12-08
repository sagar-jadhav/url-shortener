package model

type Request struct {
	LongURL string `json:"longURL"`
}

type Response struct {
	LongURL  string `json:"longURL"`
	ShortURL string `json:"shortURL"`
}
