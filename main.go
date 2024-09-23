package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"media-indexer/config"
	"media-indexer/controllers/media"
	"media-indexer/controllers/tags"
	"media-indexer/docs"
	mediaRepo "media-indexer/repositories/media"
	tagRepo "media-indexer/repositories/tag"
	mediaService "media-indexer/services/media"
	"media-indexer/services/tag"
)

func setupApp(r *gin.Engine) {
	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "healthys",
		})
	})
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	tagRepo := tagRepo.NewTagRepository(config.DB)
	mediaRepo := mediaRepo.NewMediaRepository(config.DB)

	tagService := tag.NewTagService(tagRepo)
	mediaService := mediaService.NewMediaService(mediaRepo, tagRepo)

	tagController := tags.NewTagController(tagService)
	mediaController := media.NewMediaController(mediaService)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/tags", tagController.CreateTag)
		v1.GET("/tags", tagController.ListTags)
		v1.POST("/media", mediaController.CreateMedia)
		v1.GET("/media", mediaController.SearchMediaByTag)
	}
}
func main() {
	fmt.Println("Initializing the application...")

	config.ConnectDB()
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.New()
	setupApp(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
