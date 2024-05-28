package controllers

import (
	dbase "MrScraper/db"
	"MrScraper/internal/scrapers"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

func ShowStartPage(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func UpdateNewArticles(c *gin.Context) {
	theme := c.Query("theme")

	db, err := dbase.GetInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	s := []scrapers.Scraper{
		scrapers.NewHabrScraper(),
		scrapers.NewArxScraper(),
		scrapers.NewSpringerOpenScraper(),
	}

	var wg sync.WaitGroup
	errors := make([]error, len(s))

	for i, scraper := range s {
		wg.Add(1)
		go func(i int, scraper scrapers.Scraper) {
			defer wg.Done()
			articles, err := scraper.Scrap(theme)
			if err != nil {
				errors[i] = err
				return
			}
			for _, article := range articles {
				err := db.InsertArticle(article)
				if err != nil {
					errors[i] = err
					return
				}
			}
		}(i, scraper)
	}

	wg.Wait()
	c.JSON(http.StatusOK, nil)
}

func AllArticles(c *gin.Context) {
	db, err := dbase.GetInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	articles, err := db.GetArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database get data problems"})
		return
	}

	c.JSON(http.StatusOK, articles)
}
