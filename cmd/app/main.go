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

	//router.Static("/pages", "./web/pages")
	//router.Static("/js", "./web/js")
	//router.LoadHTMLGlob("web/pages/*")

	routes.SetupRoutes(router)

	router.Run(":8080")
}
