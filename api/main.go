package main

import (
  "net/http"
  "log"

  config "konsin1988/api/config"
  health "konsin1988/api/db/health"
	http_t "konsin1988/api/transport/http"
	db "konsin1988/api/db"
)

func main() {
	config.ConnectDB()
  defer config.DB.Close()

	db.RunMigrations()

  healthRepo := health.NewHealthRepository(config.DB)
	healthService := health.NewService(healthRepo)
  healthHandler := http_t.NewHealthHandler(healthService)

  mux := http.NewServeMux()
  mux.Handle("/health", healthHandler)

  log.Println("Server started on :8013")
  log.Fatal(http.ListenAndServe("0.0.0.0:8013", mux))
}

