package request

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AydarZaynutdinov/pow/internal/logger"
)

var testLogger = logger.NewTestLogger()

func TestNewSolveRequest(t *testing.T) {
	tests := []struct {
		name    string
		reqBody string
		check   func(r *SolveRequest, err error)
	}{
		{
			name:    "success",
			reqBody: "{\"challenge\": \"testChallenge\", \"solution\": \"testSolution\"}",
			check: func(request *SolveRequest, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "testChallenge", request.Challenge)
				assert.Equal(t, "testSolution", request.Solution)
			},
		},
		{
			name:    "empty body",
			reqBody: "",
			check: func(request *SolveRequest, err error) {
				assert.Nil(t, request)
				assert.Equal(t, "failed to unmarshal request body", err.Error())
			},
		},
		{
			name:    "missing 'solution' field",
			reqBody: "\"challenge\": \"testChallenge\"",
			check: func(request *SolveRequest, err error) {
				assert.Nil(t, request)
				assert.Equal(t, "failed to unmarshal request body", err.Error())
			},
		},
		{
			name:    "missing 'challenge' field",
			reqBody: "\"solution\": \"testSolution\"",
			check: func(request *SolveRequest, err error) {
				assert.Nil(t, request)
				assert.Equal(t, "failed to unmarshal request body", err.Error())
			},
		},
		{
			name:    "empty 'solution' field",
			reqBody: "\"challenge\": \"testChallenge\", \"solution\": \"\"",
			check: func(request *SolveRequest, err error) {
				assert.Nil(t, request)
				assert.Equal(t, "failed to unmarshal request body", err.Error())
			},
		},
		{
			name:    "empty 'challenge' field",
			reqBody: "\"challenge\": \"\", \"solution\": \"testSolution\"",
			check: func(request *SolveRequest, err error) {
				assert.Nil(t, request)
				assert.Equal(t, "failed to unmarshal request body", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := prepareRequest(tt.reqBody)
			solveRequest, err := NewSolveRequest(context.Background(), req, testLogger)
			tt.check(solveRequest, err)
		})
	}
}

func prepareRequest(body string) *http.Request {
	return httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString(body))
}
