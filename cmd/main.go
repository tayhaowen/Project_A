package main

import (
	"fmt"
	"log"
	"project-a/internal/config"
	apphandler "project-a/internal/handler"
	"project-a/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

func main() {
	r := gin.Default()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	tideService := service.NewTideService(nil)
	tideHandler := apphandler.NewTideHandler(tideService)
	tideHandler.RegisterRoutes(r)

	r.GET("/check-health", func(c *gin.Context) {
		response := gin.H{
			"service": gin.H{
				"status": "ok",
				"uptime": time.Since(startTime).String(),
			},
			"timestamp": time.Now().UTC(),
		}

		c.JSON(200, response)
	})

	server := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	r.Run(server)
	log.Printf("Server is running on %s", server)
}
