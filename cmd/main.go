package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/8bye8/platform-api/db"
	"github.com/8bye8/platform-api/handlers"
	"github.com/8bye8/platform-api/middlewares"
	"github.com/8bye8/platform-api/services"
	"github.com/8bye8/platform-api/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	pubkey := os.Getenv("JWT_PUBLIC_KEY_B64")
	if pubkey == "" {
		logger.Error("JWT_PUBLIC_KEY_B64 environment variable is not set")
	}
	decodedKey, err := base64.StdEncoding.DecodeString(pubkey)
	if err != nil {
		logger.Error("failed to decode base64 pubkey", slog.Any("error", err))
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(decodedKey)
	if err != nil {
		logger.Error("failed to parse RSA public key", slog.Any("error", err))
		os.Exit(1)
	}

	pvtKeyB64 := os.Getenv("JWT_PRIVATE_KEY_B64")
	if pvtKeyB64 == "" {
		logger.Error("JWT_PRIVATE_KEY_B64 environment variable is not set")
	}
	decodedKey, err = base64.StdEncoding.DecodeString(pvtKeyB64)
	if err != nil {
		logger.Error("failed to decode base64 pvtKey", slog.Any("error", err))
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(decodedKey)
	if err != nil {
		logger.Error("failed to parse RSA private key", slog.Any("error", err))
		os.Exit(1)
	}

	databaseURL := os.Getenv("DB_CONNECTION_STRING")
	if databaseURL == "" {
		logger.Error("failed to fetch DB URL from env", slog.Any("error", err))
		os.Exit(1)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer conn.Close(ctx)

	if err := conn.Ping(ctx); err != nil {
		logger.Error("failed to ping database", slog.Any("error", err))
		os.Exit(1)
	}

	var queries db.Querier = db.New(conn)
	var hasherParams = types.HasherParams{
		SaltLength: 16,
		KeyLength:  32,
		TimeCost:   2,
		MemoryCost: 64 * 1024,
		Threads:    4,
	}

	authService := services.NewAuthService(queries, hasherParams, privateKey)
	mmService := services.NewMatchmakingService(queries)

	authHandler := handlers.NewAuthHandler(authService, logger)
	mmHandler := handlers.NewMatchmakingHandler(mmService, logger)

	router := http.NewServeMux()
	authMux := http.NewServeMux()
	mmMux := http.NewServeMux()

	authHandler.RegisterRoutes(authMux)
	mmHandler.RegisterRoutes(mmMux)

	requireAuth := middlewares.RequireAuth(publicKey)
	router.Handle("/auth/", http.StripPrefix("/auth", authMux))
	router.Handle("/mm/", http.StripPrefix("/mm", requireAuth(mmMux)))

	middlewareStack := middlewares.CreateStack(
		middlewares.Logging,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: middlewareStack(router),
	}

	logger.Info("starting server", slog.String("port", "8080"))
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error while starting server:", err)
	}
}
