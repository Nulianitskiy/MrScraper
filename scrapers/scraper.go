package scrapers

import "MrScraper/model"

type Scraper interface {
	Scrap(theme string) ([]model.Article, error)
}
