package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"pet-store/backend/internal/auth"
	"pet-store/backend/internal/config"
	"pet-store/backend/internal/db"
	"pet-store/backend/internal/graph"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gqlhandler "github.com/graphql-go/handler"
)

func runSQLFile(database *sql.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = database.Exec(string(content))
	return err
}

func ensureMigrations(database *sql.DB) error {
	files := []string{
		"internal/db/migrations/001_init.sql",
		"internal/db/migrations/002_seed.sql",
	}
	for _, f := range files {
		if err := runSQLFile(database, f); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	cfg := config.Load()

	if err := os.MkdirAll(cfg.ImageDir, 0755); err != nil {
		log.Fatal(err)
	}

	database, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := ensureMigrations(database); err != nil {
		log.Fatal(err)
	}

	schema, err := graph.NewSchema(graph.SchemaDeps{
		DB:       database,
		ImageDir: cfg.ImageDir,
	})
	if err != nil {
		log.Fatal(err)
	}

	h := gqlhandler.New(&gqlhandler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/healthz", func(c *gin.Context) {
		if err := database.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "down",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	authService := &auth.AuthService{DB: database}
	r.Use(authService.Middleware())

	r.Static("/media", filepath.Clean("./media"))

	r.GET("/query", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/query", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/playground", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/playground", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	log.Printf("server running on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
