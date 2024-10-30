package challenge

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/AydarZaynutdinov/pow/internal/config"
	"github.com/AydarZaynutdinov/pow/internal/logger"
)

var (
	testLogger = logger.NewTestLogger()
	cfg        = config.Challenge{
		Len:        1,
		TTL:        1,
		Difficulty: 1,
	}
	quotesList = []string{"1"}

	errCustom = errors.New("custom error")
)

type CacheMock struct {
	existsFunc func(ctx context.Context, key string) (int64, error)
	setFunc    func(ctx context.Context, key string, ttl time.Duration) error
	delFunc    func(ctx context.Context, key string) error
}

func (m *CacheMock) Exists(ctx context.Context, key string) (int64, error) {
	return m.existsFunc(ctx, key)
}

func (m *CacheMock) Set(ctx context.Context, key string, ttl time.Duration) error {
	return m.setFunc(ctx, key, ttl)
}

func (m *CacheMock) Del(ctx context.Context, key string) error {
	return m.delFunc(ctx, key)
}

func TestService_GenerateChallenge(t *testing.T) {
	type fields struct {
		cache Cache
	}
	tests := []struct {
		name   string
		fields fields
		check  func(result string, err error)
	}{
		{
			name: "error during saving into cache",
			fields: fields{
				cache: &CacheMock{
					setFunc: func(ctx context.Context, key string, ttl time.Duration) error {
						return errCustom
					},
				},
			},
			check: func(result string, err error) {
				assert.Empty(t, result)
				assert.Equal(t, errCustom, err)
			},
		},
		{
			name: "success",
			fields: fields{
				cache: &CacheMock{
					setFunc: func(ctx context.Context, key string, ttl time.Duration) error {
						return nil
					},
				},
			},
			check: func(result string, err error) {
				assert.NotNil(t, result)
				assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.cache, quotesList, cfg, testLogger)
			got, err := s.GenerateChallenge(context.Background())
			tt.check(got, err)
		})
	}
}

func TestService_CheckSolution(t *testing.T) {
	type fields struct {
		cache        Cache
		challengeCfg config.Challenge
	}
	type args struct {
		challenge string
		solution  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		check  func(result string, err error)
	}{
		{
			name: "error during checking challenge in cache",
			fields: fields{
				cache: &CacheMock{
					existsFunc: func(ctx context.Context, key string) (int64, error) {
						return 0, errCustom
					},
				},
				challengeCfg: config.Challenge{},
			},
			args: args{
				challenge: "",
				solution:  "",
			},
			check: func(result string, err error) {
				assert.Empty(t, result)
				assert.Equal(t, errCustom, err)
			},
		},
		{
			name: "challenge is not found",
			fields: fields{
				cache: &CacheMock{
					existsFunc: func(ctx context.Context, key string) (int64, error) {
						return 0, nil
					},
				},
				challengeCfg: config.Challenge{},
			},
			args: args{
				challenge: "",
				solution:  "",
			},
			check: func(result string, err error) {
				assert.Empty(t, result)
				assert.Equal(t, ErrChallengeNotFound, err)
			},
		},
		{
			name: "incorrect solution for challenge",
			fields: fields{
				cache: &CacheMock{
					existsFunc: func(ctx context.Context, key string) (int64, error) {
						return 1, nil
					},
				},
				challengeCfg: config.Challenge{
					Difficulty: 10,
				},
			},
			args: args{
				challenge: "1",
				solution:  "2",
			},
			check: func(result string, err error) {
				assert.Empty(t, result)
				assert.Equal(t, ErrIncorrectSolution, err)
			},
		},
		{
			name: "correct solution",
			fields: fields{
				cache: &CacheMock{
					existsFunc: func(ctx context.Context, key string) (int64, error) {
						return 1, nil
					},
					delFunc: func(ctx context.Context, key string) error {
						return nil
					},
				},
				challengeCfg: config.Challenge{
					Difficulty: 0,
				},
			},
			args: args{
				challenge: "1",
				solution:  "2",
			},
			check: func(result string, err error) {
				assert.Equal(t, quotesList[0], result)
				assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.cache, quotesList, tt.fields.challengeCfg, testLogger)
			got, err := s.CheckSolution(context.Background(), tt.args.challenge, tt.args.solution)
			tt.check(got, err)
		})
	}
}

func TestService_IsSolutionIncorrect(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		check func(result bool)
	}{
		{
			name: "ErrIncorrectSolution",
			args: args{
				err: ErrIncorrectSolution,
			},
			check: func(result bool) {
				assert.True(t, result)
			},
		},
		{
			name: "ErrChallengeNotFound",
			args: args{
				err: ErrChallengeNotFound,
			},
			check: func(result bool) {
				assert.True(t, result)
			},
		},
		{
			name: "nil",
			args: args{
				err: nil,
			},
			check: func(result bool) {
				assert.False(t, result)
			},
		},
		{
			name: "another error",
			args: args{
				err: errors.New("another error"),
			},
			check: func(result bool) {
				assert.False(t, result)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(nil, quotesList, cfg, testLogger)
			got := s.IsSolutionIncorrect(tt.args.err)
			tt.check(got)
		})
	}
}

func TestService_GetDifficulty(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.Challenge
		want int
	}{
		{
			name: "0",
			cfg: config.Challenge{
				Difficulty: 0,
			},
			want: 0,
		},
		{
			name: "1",
			cfg: config.Challenge{
				Difficulty: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(nil, quotesList, tt.cfg, testLogger)
			got := s.GetDifficulty()
			assert.Equal(t, tt.want, got)
		})
	}
}
