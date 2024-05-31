package main

import (
	dbase "MrScraper/db"
	"MrScraper/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	db, err := dbase.GetInstance()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := gin.Default()

	router.LoadHTMLGlob("web/pages/*")

	router.Static("/pages", "./web/pages")
	router.Static("/js", "./web/js")
	router.Static("/styles", "./web/styles")

	routes.SetupRoutes(router)

	router.Run(":8080")
}
