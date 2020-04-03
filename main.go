package main

import (
	"covid19/external"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	engine := gin.Default()
	engine.LoadHTMLGlob("templates/*")
	engine.Static("/static","./static")
	engine.GET("/", func(context *gin.Context) {
		coronaUpdate := external.LiveData()
		context.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"index.html",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"title": "COVID-19",
				"data" : coronaUpdate,
			},
		)
	})
	engine.GET("/crowdsourced", func(context *gin.Context) {
		coronaUpdate := external.CrowdData()
		context.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"crowdsource.html",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"title": "COVID-19",
				"data" : coronaUpdate,
			},
		)
	})

	_ = engine.Run()
}
