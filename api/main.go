package main

import (
  "net/http"
  "log"

  config "konsin1988/api/config"
  health "konsin1988/api/db/health"
	auth "konsin1988/api/auth"
	http_t "konsin1988/api/transport/http"
	db "konsin1988/api/db"
)

func main() {
	config.ConnectDB()
  defer config.DB.Close()

	db.RunMigrations()

	err := auth.InitKeycloak()
	if err != nil {
		log.Fatal(err)
	}

  healthRepo := health.NewHealthRepository(config.DB)
	healthService := health.NewService(healthRepo)
  healthHandler := http_t.NewHealthHandler(healthService)

	api := http.NewServeMux()
	
	//api.HandleFunc("GET /containers", containerHandler.ListContainers)
	//api.Handle("/nodes", nodeHandler)
	
	
	api.Handle("/health", healthHandler)
	//root.Handle("/api/", auth.Middleware(api))

  log.Println("Server started on :8013")
	log.Fatal(http.ListenAndServe("0.0.0.0:8013", api))
}

