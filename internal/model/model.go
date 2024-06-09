package model

type Article struct {
	Id         int    `db:"id" json:"id"`
	Title      string `db:"title" json:"title"`
	Authors    string `db:"author" json:"authors"`
	Annotation string `db:"abstract" json:"annotation"`
	Text       string `db:"content" json:"text,omitempty"`
	Link       string `db:"url" json:"link"`
	Theme      string `db:"theme" json:"theme,omitempty"`
}
