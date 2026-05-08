package http 

import (
    "encoding/json"
    "net/http"
    "context"
)

type HealthService interface {
  Check(ctx context.Context) error
}

type HealthHandler struct {
  service HealthService
}

func NewHealthHandler(service HealthService) *HealthHandler {
  return &HealthHandler{service: service}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  err := h.service.Check(r.Context())

  w.Header().Set("Content-Type", "application/json")

  if err != nil {
    w.WriteHeader(http.StatusServiceUnavailable)
    json.NewEncoder(w).Encode(map[string]string{
      "status": "down",
    })

    return 
  }
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(map[string]string{
    "status": "ok",
  })
}


