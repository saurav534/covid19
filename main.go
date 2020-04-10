package main

import (
	"context"
	"covid19/external"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cron *external.Cron
var srv *http.Server

func main() {
	engine := gin.Default()
	engine.LoadHTMLGlob("templates/*")
	engine.Static("/static", "./static")
	engine.RedirectFixedPath = true
	external.UpdateMohData()
	external.UpdateCrowdData()
	engine.GET("/", func(context *gin.Context) {
		coronaUpdate := external.GetCrowdData()
		context.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the moh.html template
			"crowdsource.html",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"title":   "COVID-19",
				"nextTab": "/moh",
				"data":    coronaUpdate,
			},
		)
	})
	engine.GET("/crowdsourced", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/")
		context.Abort()
	})
	engine.GET("/moh", func(context *gin.Context) {
		coronaUpdate := external.GetMohData()
		context.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the moh.html template
			"moh.html",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"title":   "COVID-19",
				"nextTab": "/",
				"data":    coronaUpdate,
			},
		)
	})

	cron = external.NewCron()
	go func() {
		cron.StartCron()
	}()
	stopSignal := watchStopSignal()

	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// starting server
	srv = &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	log.Println("Server stopped successfully")
	<-stopSignal
}

func watchStopSignal() chan bool {
	stopSignal := make(chan os.Signal)
	doneClosure := make(chan bool)

	signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-stopSignal
		log.Printf("received signal %s, waiting to drain requests", sig)
		cron.StopCron()

		// stopping the server
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		doneClosure <- true
	}()
	return doneClosure
}
