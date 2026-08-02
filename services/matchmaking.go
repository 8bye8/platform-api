package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/8bye8/platform-api/db"
	"github.com/8bye8/platform-api/middlewares"
	"github.com/8bye8/platform-api/types"
	"github.com/jackc/pgx/v5/pgtype"
)

type MatchmakingService struct {
	db db.Querier
}

func NewMatchmakingService(database db.Querier) *MatchmakingService {
	return &MatchmakingService{db: database}
}

func (s *MatchmakingService) CreateChallenge(ctx context.Context, params types.ChallengeParams) error {
	userId, ok := ctx.Value(middlewares.AuthUserId).(string)
	if !ok {
		slog.Error("UserId not found in context")
		return errors.New("UserId not found in context")
	}
	var createdBy pgtype.UUID
	if err := createdBy.Scan(userId); err != nil {
		return fmt.Errorf("invalid authenticated user ID: %w", err)
	}

	challengeParams := db.CreateChallengeParams{
		CreatedBy:   createdBy,
		MinRating:   int32(params.MinRating),
		MaxRating:   int32(params.MaxRating),
		Rating:      int32(params.Rating),
		SideToPlay:  db.ChallengeSideToPlay(params.SideToPlay),
		TimeControl: db.ChallengeTimeControl(params.TimeControl),
	}
	_, err := s.db.CreateChallenge(ctx, challengeParams)
	if err != nil {
		return err
	}
	return nil
}
