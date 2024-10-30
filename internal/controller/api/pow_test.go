package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AydarZaynutdinov/pow/internal/controller/api/response"

	"github.com/stretchr/testify/assert"

	"github.com/AydarZaynutdinov/pow/internal/controller/api/utils"
	"github.com/AydarZaynutdinov/pow/internal/logger"
)

var (
	testLogger = logger.NewTestLogger()

	errCustom = errors.New("custom error")
)

type ServiceMock struct {
	generateChallengeFunc   func(ctx context.Context) (string, error)
	checkSolutionFunc       func(ctx context.Context, challenge, solution string) (string, error)
	isSolutionIncorrectFunc func(err error) bool
	getDifficultyFunc       func() int
}

func (m *ServiceMock) GenerateChallenge(ctx context.Context) (string, error) {
	return m.generateChallengeFunc(ctx)
}

func (m *ServiceMock) CheckSolution(ctx context.Context, challenge, solution string) (string, error) {
	return m.checkSolutionFunc(ctx, challenge, solution)
}

func (m *ServiceMock) IsSolutionIncorrect(err error) bool {
	return m.isSolutionIncorrectFunc(err)
}

func (m *ServiceMock) GetDifficulty() int {
	return m.getDifficultyFunc()
}

func TestController_GetChallenge(t *testing.T) {
	tests := []struct {
		name       string
		service    Service
		prepareReq func() *http.Request
		check      func(result *utils.ControllerResult, err error)
	}{
		{
			name: "error during generating challenge",
			service: &ServiceMock{
				generateChallengeFunc: func(ctx context.Context) (string, error) {
					return "", errCustom
				},
			},
			prepareReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/challenge", nil)
			},
			check: func(result *utils.ControllerResult, err error) {
				assert.Equal(t, utils.NewControllerStandardErrorInternalServerError(errCustom), err)
				assert.Nil(t, result)
			},
		},
		{
			name: "success",
			service: &ServiceMock{
				generateChallengeFunc: func(ctx context.Context) (string, error) {
					return "challenge", nil
				},
				getDifficultyFunc: func() int {
					return 1
				},
			},
			prepareReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/challenge", nil)
			},
			check: func(result *utils.ControllerResult, err error) {
				assert.Nil(t, err)
				assert.Equal(t, utils.NewControllerResultOK(response.ChallengeResponse{
					Challenge:  "challenge",
					Difficulty: 1,
				}), result)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				service: tt.service,
				logger:  testLogger,
			}
			req := tt.prepareReq()
			got, err := c.GetChallenge(req)
			tt.check(got, err)
		})
	}
}

func TestController_Solve(t *testing.T) {
	tests := []struct {
		name       string
		service    Service
		prepareReq func() *http.Request
		check      func(result *utils.ControllerResult, err error)
	}{
		{
			name:    "error during preparing solve request",
			service: &ServiceMock{},
			prepareReq: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/solve", http.NoBody)
			},
			check: func(result *utils.ControllerResult, err error) {
				assert.Equal(t, utils.NewControllerStandardErrorBadRequest(fmt.Errorf("failed to unmarshal request body")), err)
				assert.Nil(t, result)
			},
		},
		{
			name: "error during preparing solve request | incorrect solution",
			service: &ServiceMock{
				checkSolutionFunc: func(ctx context.Context, challenge, solution string) (string, error) {
					return "", errCustom
				},
				isSolutionIncorrectFunc: func(err error) bool {
					return true
				},
			},
			prepareReq: func() *http.Request {
				payload := "{\"challenge\": \"testChallenge\", \"solution\": \"testSolution\"}"
				return httptest.NewRequest(http.MethodPost, "/solve", bytes.NewBuffer([]byte(payload)))
			},
			check: func(result *utils.ControllerResult, err error) {
				assert.Equal(t, utils.NewControllerStandardErrorBadRequest(errCustom), err)
				assert.Nil(t, result)
			},
		},
		{
			name: "error during preparing solve request | another reason",
			service: &ServiceMock{
				checkSolutionFunc: func(ctx context.Context, challenge, solution string) (string, error) {
					return "", errCustom
				},
				isSolutionIncorrectFunc: func(err error) bool {
					return false
				},
			},
			prepareReq: func() *http.Request {
				payload := "{\"challenge\": \"testChallenge\", \"solution\": \"testSolution\"}"
				return httptest.NewRequest(http.MethodPost, "/solve", bytes.NewBuffer([]byte(payload)))
			},
			check: func(result *utils.ControllerResult, err error) {
				assert.Equal(t, utils.NewControllerStandardErrorInternalServerError(errCustom), err)
				assert.Nil(t, result)
			},
		},
		{
			name: "success",
			service: &ServiceMock{
				checkSolutionFunc: func(ctx context.Context, challenge, solution string) (string, error) {
					return "quote 1", nil
				},
			},
			prepareReq: func() *http.Request {
				payload := "{\"challenge\": \"testChallenge\", \"solution\": \"testSolution\"}"
				return httptest.NewRequest(http.MethodPost, "/solve", bytes.NewBuffer([]byte(payload)))
			},
			check: func(result *utils.ControllerResult, err error) {
				assert.Equal(t, utils.NewControllerResultOK(
					response.SolveResponse{
						Quote: "quote 1",
					},
				), result)
				assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				service: tt.service,
				logger:  testLogger,
			}
			req := tt.prepareReq()
			got, err := c.Solve(req)
			tt.check(got, err)
		})
	}
}
