package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"errors"
	"fmt"
	"time"

	"github.com/8bye8/platform-api/db"
	"github.com/8bye8/platform-api/types"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

type AuthService struct {
	db           db.Querier
	hasherParams types.HasherParams
	jwtPvtKey    *rsa.PrivateKey
}

func NewAuthService(database db.Querier, params types.HasherParams, jwtPvtKey *rsa.PrivateKey) *AuthService {
	return &AuthService{
		db:           database,
		hasherParams: params,
		jwtPvtKey:    jwtPvtKey,
	}
}

func (s *AuthService) Register(ctx context.Context, params types.RegisterUserParams) (string, error) {
	salt := make([]byte, s.hasherParams.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	pwdHash := argon2.IDKey(
		[]byte(params.Password),
		salt,
		s.hasherParams.TimeCost,
		s.hasherParams.MemoryCost,
		s.hasherParams.Threads,
		s.hasherParams.KeyLength,
	)
	userParams := db.CreateUserParams{
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: pwdHash,
		PasswordSalt: salt,
	}
	user, err := s.db.CreateUser(ctx, userParams)
	if err != nil {
		return "", err
	}
	return user.ID.String(), err
}

func (s *AuthService) Login(ctx context.Context, params types.LoginParams) (string, error) {
	user, err := s.db.GetUserByUsername(ctx, params.Username)
	if err != nil {
		return "", err
	}

	newHash := argon2.IDKey(
		[]byte(params.Password),
		user.PasswordSalt,
		s.hasherParams.TimeCost,
		s.hasherParams.MemoryCost,
		s.hasherParams.Threads,
		s.hasherParams.KeyLength,
	)

	match := subtle.ConstantTimeCompare(newHash, user.PasswordHash)
	if match != 1 {
		return "", errors.New("invalid password")
	}

	claims := jwt.MapClaims{
		"sub": user.ID.String(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(s.jwtPvtKey)
	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
		return "", err
	}

	return tokenString, nil
}
