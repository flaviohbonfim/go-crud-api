package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"go-crud-api/internal/config"
	"go-crud-api/internal/database"
	"go-crud-api/internal/domain/products"
	"go-crud-api/internal/domain/users"
	customhttp "go-crud-api/internal/http"
	"go-crud-api/internal/logger"
	"go-crud-api/internal/repository"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a JWT.
func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load environment variables")
	}

	// Init Logger
	logger.InitLogger(cfg.AppEnv)

	log.Info().Msg("Starting application...")

	// Connect to the database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to the database")
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	log.Info().Msg("Database connection successful.")

	// Run Migrations
	runMigrations(sqlDB)

	// Dependency Injection
	userRepo := repository.NewGormUserRepository(db)
	userService := users.NewService(userRepo, cfg)
	authHandler := users.NewAuthHandler(userService)

	productRepo := repository.NewGormProductRepository(db)
	productService := products.NewService(productRepo)
	productHandler := products.NewProductHandler(productService)

	// Initialize Router
	router := customhttp.InitRouter(cfg, db, authHandler, productHandler)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Info().Msgf("Server starting on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}

func runMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create migrate driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create migrate instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("Could not apply migrations")
	}

	log.Info().Msg("Database migrations applied successfully.")
}
