package main

import (
	"log"
	"net/http"
	"net/http/pprof"
	"os"

	auth "konsin1988/api/auth"
	config "konsin1988/api/config"
	db "konsin1988/api/db"
	health "konsin1988/api/db/health"
	http_t "konsin1988/api/transport/http"
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

	proxyPath := os.Getenv("PROXY_PATH")

	api.Handle(proxyPath+"/health", healthHandler)
	//root.Handle("/api/", auth.Middleware(api))
	// Register pprof endpoints.

	api.HandleFunc("/debug/pprof/", pprof.Index)
	api.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	api.HandleFunc("/debug/pprof/profile", pprof.Profile)
	api.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	api.HandleFunc("/debug/pprof/trace", pprof.Trace)

	api.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	api.Handle("/debug/pprof/block", pprof.Handler("block"))
	api.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	api.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	api.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	api.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	log.Println("Server started on :8013")
	log.Fatal(http.ListenAndServe("0.0.0.0:8013", api))
}
