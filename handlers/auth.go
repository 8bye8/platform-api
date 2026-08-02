package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/8bye8/platform-api/types"
)

type AuthService interface {
	Register(ctx context.Context, params types.RegisterUserParams) (string, error)
	Login(ctx context.Context, params types.LoginParams) (string, error)
}

type AuthHandler struct {
	svc    AuthService
	logger *slog.Logger
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func NewAuthHandler(svc AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		svc:    svc,
		logger: logger,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "handling user registration",
		slog.String("path", r.URL.Path),
		slog.String("method", r.Method),
	)
	var registerReq RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	userId, err := h.svc.Register(r.Context(), types.RegisterUserParams{
		Username: registerReq.Username,
		Password: registerReq.Password,
		Email:    registerReq.Email,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(userId))
	if err != nil {
		return
	}
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "handling user login")
	// Decode directly from the stream into the struct
	var loginReq LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	token, err := h.svc.Login(r.Context(), types.LoginParams{
		Username: loginReq.Username,
		Password: loginReq.Password,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(token))
	if err != nil {
		return
	}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", h.RegisterUser)
	mux.HandleFunc("POST /login", h.LoginUser)
}
