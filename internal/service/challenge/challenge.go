package challenge

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/AydarZaynutdinov/pow/internal/config"
	"github.com/AydarZaynutdinov/pow/internal/logger"

	"go.uber.org/zap"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	ErrIncorrectSolution = errors.New("incorrect solution")
	ErrChallengeNotFound = errors.New("challenge not found")
)

type Cache interface {
	Exists(ctx context.Context, key string) (int64, error)
	Set(ctx context.Context, key string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

type Service struct {
	cache        Cache
	quotesList   []string
	challengeCfg config.Challenge
	logger       logger.Logger
}

func NewService(
	cache Cache,
	quotesList []string,
	challengeCfg config.Challenge,
	logger logger.Logger,
) *Service {
	return &Service{
		cache:        cache,
		quotesList:   quotesList,
		challengeCfg: challengeCfg,
		logger:       logger,
	}
}

func (s *Service) GenerateChallenge(ctx context.Context) (string, error) {
	challenge := s.generateChallenge()

	err := s.cache.Set(ctx, challenge, s.challengeCfg.TTL)
	if err != nil {
		s.logger.WithContext(ctx).Error(
			"failed to set challenge into cache",
			zap.Error(err),
		)
		return "", err
	}

	return challenge, nil
}

func (s *Service) CheckSolution(ctx context.Context, challenge, solution string) (string, error) {
	exists, err := s.cache.Exists(ctx, challenge)
	if err != nil {
		s.logger.WithContext(ctx).Error(
			"failed to check challenge in cache",
			zap.Error(err),
		)
		return "", err
	}
	if exists == 0 {
		return "", ErrChallengeNotFound
	}

	if !s.checkSolution(challenge, solution) {
		return "", ErrIncorrectSolution
	}

	err = s.cache.Del(ctx, challenge)
	if err != nil {
		s.logger.WithContext(ctx).Warn(
			"failed to delete challenge from cache",
			zap.Error(err),
		)
	}

	return s.getRandomQuote(), nil
}

func (s *Service) IsSolutionIncorrect(err error) bool {
	return errors.Is(err, ErrIncorrectSolution) || errors.Is(err, ErrChallengeNotFound)
}

func (s *Service) GetDifficulty() int {
	return s.challengeCfg.Difficulty
}

func (s *Service) generateChallenge() string {
	var challenge strings.Builder
	for i := 0; i < s.challengeCfg.Len; i++ {
		ind, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		challenge.WriteByte(chars[ind.Int64()])
	}
	return challenge.String()
}

func (s *Service) checkSolution(challenge, solution string) bool {
	hash := sha256.Sum256([]byte(challenge + solution))
	hashStr := hex.EncodeToString(hash[:])
	return strings.HasPrefix(hashStr, strings.Repeat("0", s.challengeCfg.Difficulty))
}

func (s *Service) getRandomQuote() string {
	ind, _ := rand.Int(rand.Reader, big.NewInt(int64(len(s.quotesList))))
	return s.quotesList[ind.Int64()]
}
