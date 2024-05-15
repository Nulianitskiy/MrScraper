package scrapers

import (
	"MrScraper/model"
	"MrScraper/utils"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

type ArxivScraper struct {
}

func NewArxScraper() *ArxivScraper {
	return &ArxivScraper{}
}

func (s ArxivScraper) Scrap(theme string) ([]model.Article, error) {
	var articles []model.Article

	c := colly.NewCollector(
		colly.AllowedDomains("arxiv.org"),
	)

	// Обработка каждого элемента статьи на странице поиска
	c.OnHTML(".arxiv-result", func(e *colly.HTMLElement) {
		// Извлекаем заголовок
		title := strings.TrimSpace(e.ChildText(".title"))
		// Извлекаем авторов
		authors := strings.TrimSpace(e.ChildText(".authors"))
		authors = strings.Replace(authors, "Authors:", "", -1)
		authors = strings.Replace(authors, "\n", "", -1)
		// Извлекаем ссылку на статью
		articleLink := e.ChildAttr(".list-title.is-inline-block a", "href")

		pdfLink := e.ChildAttr(".list-title.is-inline-block span a", "href")
		fmt.Println(pdfLink)

		// Извлечение текста из PDF
		res, err := utils.ExtractTextFromPDF(pdfLink)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		// Добавляем промежуточную структуру с заголовком и авторами в слайс
		articles = append(articles, model.Article{
			Title:   title,
			Authors: authors,
			Link:    articleLink,
			Text:    res,
		})

		// Переходим на страницу статьи для получения дополнительной информации
		e.Request.Visit(articleLink)
	})

	// Обработка самой статьи на странице статьи
	c.OnHTML("div#content", func(e *colly.HTMLElement) {
		// Извлекаем аннотацию
		abstract := strings.TrimSpace(e.ChildText("blockquote.abstract"))
		abstract = strings.Replace(abstract, "Abstract:", "", -1)

		// Найти соответствующую статью по ссылке и обновить её
		for i := range articles {
			if articles[i].Link == e.Request.URL.String() {
				articles[i].Annotation = abstract
				break
			}
		}
	})

	// Обработка ошибок
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Формируем URL для поиска статей по теме
	searchURL := fmt.Sprintf("https://arxiv.org/search/?query=%s&searchtype=all&abstracts=show&order=-announced_date_first&size=25", theme)

	//c.OnScraped(func(r *colly.Response) {
	//	for _, article := range articles {
	//		fmt.Printf("Title: %s\n", article.Title)
	//		fmt.Printf("Authors: %s\n", article.Authors)
	//		fmt.Printf("Annotation: %s\n", article.Annotation)
	//		fmt.Printf("Link: %s\n", article.Link)
	//		fmt.Printf("Text: %s\n\n", article.Text)
	//	}
	//})

	// Запускаем сбор данных
	err := c.Visit(searchURL)
	if err != nil {
		log.Printf("Error scraping URL %s: %v", searchURL, err)
		return articles, err
	}

	return articles, nil
}
