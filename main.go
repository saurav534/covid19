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
	engine.RedirectFixedPath = true
	engine.GET("/", func(context *gin.Context) {
		coronaUpdate := external.CrowdData()
		context.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the moh.html template
			"crowdsource.html",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"title": "COVID-19",
				"nextTab" : "/moh",
				"data" : coronaUpdate,
			},
		)
	})
	engine.GET("/crowdsourced", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/")
		context.Abort()
	})
	engine.GET("/moh", func(context *gin.Context) {
		coronaUpdate := external.LiveData()
		context.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the moh.html template
			"moh.html",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"title": "COVID-19",
				"nextTab" : "/",
				"data" : coronaUpdate,
			},
		)
	})

	_ = engine.Run()
}
