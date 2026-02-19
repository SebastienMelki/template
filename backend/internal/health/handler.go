package health

import (
	"encoding/json"
	"net/http"
)

// RegisterRoutes adds health check endpoints to the given mux.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", handleLiveness)
	mux.HandleFunc("GET /readyz", handleReadiness)
}

func handleLiveness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleReadiness(w http.ResponseWriter, _ *http.Request) {
	// TODO: Add actual readiness checks (database connectivity, etc.)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}
