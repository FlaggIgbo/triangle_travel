package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"triangle_travel/internal/api"
	"triangle_travel/internal/db"
)

// Run starts the HTTP server
func Run() {
	host := flag.String("host", "localhost", "Server host")
	port := flag.Int("port", 8080, "Server port")
	dataDir := flag.String("data", ".", "Project root (contains db/data.sqlite3)")
	flag.Parse()

	// Render.com and other PaaS set PORT
	if p := os.Getenv("PORT"); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			*port = n
		}
	}

	database, err := db.New(*dataDir)
	if err != nil {
		log.Fatalf("Database: %v", err)
	}
	defer database.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	handlers := &api.Handlers{DB: database}
	apiGroup := router.Group("/api")
	apiGroup.POST("/search", handlers.Search)
	apiGroup.GET("/cities", handlers.Cities)

	buildPath := filepath.Join(*dataDir, "build")
	if _, err := os.Stat(buildPath); err == nil {
		router.Use(static.Serve("/", static.LocalFile(buildPath, true)))
		router.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(buildPath, "index.html"))
		})
	}

	serverPath := fmt.Sprintf("%s:%d", *host, *port)
	server := &http.Server{
		Addr:         serverPath,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Server started at http://%s\n", serverPath)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
	log.Println("Server exiting")
}
