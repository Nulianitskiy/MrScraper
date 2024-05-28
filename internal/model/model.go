package model

type Article struct {
	Title      string `json:"title"`
	Authors    string `json:"authors"`
	Annotation string `json:"annotation"`
	Text       string `json:"text"`
	Link       string `json:"link"`
}
