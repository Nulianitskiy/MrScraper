package scrapers

import (
	"MrScraper/internal/model"
	"MrScraper/internal/utils"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

type SpringerOpenScraper struct{}

func NewSpringerOpenScraper() *SpringerOpenScraper { return &SpringerOpenScraper{} }

func (s SpringerOpenScraper) Scrap(theme string) ([]model.Article, error) {
	var articles []model.Article

	c := colly.NewCollector(
		colly.AllowedDomains("www.springeropen.com"),
	)

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36",
	}
	uaIndex := 0

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgents[uaIndex])
		// Увеличиваем индекс для следующего запроса
		uaIndex = (uaIndex + 1) % len(userAgents)
	})

	c.OnHTML(".c-listing__content.u-mb-16", func(e *colly.HTMLElement) {
		article := model.Article{}

		article.Title = strings.TrimSpace(e.ChildText("a[data-test=title-link]"))
		article.Link = "https:" + e.ChildAttr("a[data-test=title-link]", "href")
		article.Annotation = strings.TrimSpace(e.ChildText("p"))
		article.Authors = strings.TrimSpace(e.ChildText(".c-listing__authors-list"))

		//мне лень вытаскивать ссылку снова
		pdfLink := strings.Replace(article.Link, "articles", "counter/pdf", -1) + ".pdf"
		fmt.Println(pdfLink)
		res, err := utils.ExtractTextFromPDF(pdfLink)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}
		article.Text = res

		articles = append(articles, article)
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

	searchURL := fmt.Sprintf("https://www.springeropen.com/search?query=%s&searchType=publisherSearch", theme)
	// Start scraping the specified URL
	err := c.Visit(searchURL)
	if err != nil {
		log.Fatal(err)
	}

	return articles, nil
}
