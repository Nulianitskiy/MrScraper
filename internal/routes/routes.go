package routes

import (
	"MrScraper/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", controllers.ShowStartPage)
	r.GET("/articles", controllers.AllArticles)
	r.GET("/article")

	r.GET("/update", controllers.UpdateNewArticles)
}
