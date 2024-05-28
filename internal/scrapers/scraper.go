package scrapers

import (
	"MrScraper/internal/model"
)

type Scraper interface {
	Scrap(theme string) ([]model.Article, error)
}
