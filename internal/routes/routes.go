package routes

import (
	"MrScraper/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", controllers.ShowStartPage)
	r.GET("/article", controllers.ShowArticle)
	r.GET("/check", controllers.ShowCheck)

	r.GET("/articles", controllers.AllArticles)
	r.GET("/articles/:id", controllers.ArticleById)

	r.POST("/update", controllers.UpdateNewArticles)
}
