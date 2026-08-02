package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/8bye8/platform-api/types"
)

type MatchmakingService interface {
	CreateChallenge(ctx context.Context, conditions types.ChallengeParams) error
}

type MatchmakingHandler struct {
	svc    MatchmakingService
	logger *slog.Logger
}

func NewMatchmakingHandler(svc MatchmakingService, logger *slog.Logger) *MatchmakingHandler {
	return &MatchmakingHandler{
		svc:    svc,
		logger: logger,
	}
}

func (h *MatchmakingHandler) CreateChallenge(w http.ResponseWriter, r *http.Request) {

	var params types.ChallengeParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	if err := params.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.CreateChallenge(r.Context(), params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MatchmakingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /challenge", h.CreateChallenge)
}
