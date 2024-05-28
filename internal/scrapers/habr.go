package scrapers

import (
	"MrScraper/internal/model"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

type HabrScraper struct {
}

func NewHabrScraper() *HabrScraper { return &HabrScraper{} }

func (s HabrScraper) Scrap(theme string) ([]model.Article, error) {
	var articles []model.Article

	c := colly.NewCollector(
		colly.AllowedDomains("habr.com"),
		colly.UserAgent("user-agent':'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"),
	)

	c.OnHTML(".tm-articles-list__item", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("a.tm-title__link"))
		link := fmt.Sprintf("https://habr.com%s", e.ChildAttr("a.tm-title__link", "href"))
		authors := strings.TrimSpace(e.ChildText("a.tm-user-info__username"))
		annotation := strings.TrimSpace(e.ChildText(".article-formatted-body"))
		annotation = strings.Replace(annotation, "\n", "", -1)

		articles = append(articles, model.Article{
			Title:      title,
			Authors:    authors,
			Link:       link,
			Annotation: annotation,
		})

		e.Request.Visit(link)
	})

	c.OnHTML(".tm-article-body", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.ChildText(".article-formatted-body"))

		for i := range articles {
			if articles[i].Link == e.Request.URL.String() {
				articles[i].Text = text
				break
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	//c.OnScraped(func(r *colly.Response) {
	//	for _, article := range articles {
	//		fmt.Printf("Title: %s\n", article.Title)
	//		fmt.Printf("Authors: %s\n", article.Authors)
	//		fmt.Printf("Annotation: %s\n", article.Annotation)
	//		fmt.Printf("Link: %s\n", article.Link)
	//		fmt.Printf("Text: %s\n\n", article.Text)
	//	}
	//})

	searchURL := fmt.Sprintf("https://habr.com/ru/search/?q=%s&target_type=posts&order=relevance", theme)

	err := c.Visit(searchURL)
	if err != nil {
		log.Printf("Error scraping URL %s: %v", searchURL, err)
		return articles, err
	}

	return articles, nil
}
