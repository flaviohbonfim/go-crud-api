package http

import (
	"net/http"

	"go-crud-api/pkg/web"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// HealthCheckHandler returns the health status of the application.
func HealthCheckHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ping the database
		sqlDB, err := db.DB()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get underlying DB from GORM")
			web.RespondWithError(w, "db_error", "Database connection error", http.StatusInternalServerError)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			log.Error().Err(err).Msg("Failed to ping database")
			web.RespondWithError(w, "db_error", "Database connection error", http.StatusInternalServerError)
			return
		}

		web.RespondWithJSON(w, http.StatusOK, web.Response{Data: map[string]string{"status": "ok", "db_status": "ok"}})
	}
}
