package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/service"
)

// Handler contains HTTP handlers for the example module.
// TODO: Replace with sebuf-generated service interface when contracts are generated.
// The generated code will provide a RegisterExampleServiceServer function and an
// ExampleServiceServer interface that this handler should implement.
type Handler struct {
	service *service.ExampleService
}

// NewHandler creates a new HTTP handler.
func NewHandler(svc *service.ExampleService) *Handler {
	return &Handler{service: svc}
}

// RegisterRoutes registers this module's HTTP routes.
// TODO: Replace with sebuf-generated route registration:
//
//	examplev1.RegisterExampleServiceServer(handler, examplev1.WithMux(mux))
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/example/v1/create-example", h.handleCreate)
	mux.HandleFunc("POST /api/example/v1/get-example", h.handleGetByID)
	mux.HandleFunc("POST /api/example/v1/list-examples", h.handleList)
	mux.HandleFunc("POST /api/example/v1/delete-example", h.handleDelete)
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	example, err := h.service.Create(r.Context(), req.Name, req.Description)
	if err != nil {
		slog.Error("failed to create example", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"example_id": map[string]string{"value": example.ID.String()},
	})
}

func (h *Handler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ExampleID struct {
			Value string `json:"value"`
		} `json:"example_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.ExampleID.Value)
	if err != nil {
		http.Error(w, "invalid example_id", http.StatusBadRequest)
		return
	}

	example, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("failed to get example", "id", id, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if example == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"example": map[string]any{
			"id":          map[string]string{"value": example.ID.String()},
			"name":        example.Name,
			"description": example.Description,
			"created_at":  example.CreatedAt,
			"updated_at":  example.UpdatedAt,
		},
	})
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Limit  int32 `json:"limit"`
		Offset int32 `json:"offset"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	examples, total, err := h.service.List(r.Context(), req.Limit, req.Offset)
	if err != nil {
		slog.Error("failed to list examples", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	items := make([]map[string]any, 0, len(examples))
	for _, e := range examples {
		items = append(items, map[string]any{
			"id":          map[string]string{"value": e.ID.String()},
			"name":        e.Name,
			"description": e.Description,
			"created_at":  e.CreatedAt,
			"updated_at":  e.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"examples":    items,
		"total_count": total,
	})
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ExampleID struct {
			Value string `json:"value"`
		} `json:"example_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.ExampleID.Value)
	if err != nil {
		http.Error(w, "invalid example_id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		slog.Error("failed to delete example", "id", id, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{})
}
