package test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/AydarZaynutdinov/pow/internal/controller/api"

	"github.com/AydarZaynutdinov/pow/internal/controller/api/response"
	jsoniter "github.com/json-iterator/go"

	"github.com/AydarZaynutdinov/pow/internal/config"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

const (
	defaultConfigPath = "./config/config.yaml"
)

var (
	baseURL string
	cfg     *config.Config

	httpClient *http.Client

	validate = validator.New()
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	var err error
	cfg, err = config.New(defaultConfigPath, validate)
	require.NoError(t, err)

	baseURL = fmt.Sprintf("http://%s:%d/api/pow", cfg.AppServer.Host, cfg.AppServer.Port)

	httpClient = http.DefaultClient

	testRequestChallenge(t)
	testSolveChallenge(t)
}

func testRequestChallenge(t *testing.T) {
	testCases := []struct {
		name  string
		check func(resp *http.Response)
	}{
		{
			name: "success",
			check: func(resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
				defer func() {
					_ = resp.Body.Close()
				}()

				payload, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var realResp response.ChallengeResponse
				err = jsoniter.Unmarshal(payload, &realResp)
				require.NoError(t, err)

				require.NotEmpty(t, realResp.Challenge)
				require.Equal(t, cfg.Challenge.Len, len(realResp.Challenge))
				require.Equal(t, cfg.Challenge.Difficulty, realResp.Difficulty)
			},
		},
	}

	url := fmt.Sprintf("%s/challenge", baseURL)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := httpClient.Get(url)
			require.NoError(t, err)

			tc.check(resp)
		})
	}
}

func testSolveChallenge(t *testing.T) {
	getChallenge := func() (string, int) {
		url := fmt.Sprintf("%s/challenge", baseURL)
		resp, err := httpClient.Get(url)
		require.NoError(t, err)

		defer func() {
			_ = resp.Body.Close()
		}()

		payload, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var realResp response.ChallengeResponse
		err = jsoniter.Unmarshal(payload, &realResp)
		require.NoError(t, err)

		return realResp.Challenge, realResp.Difficulty
	}

	testCases := []struct {
		name  string
		solve func(string, int) (string, string)
		check func(resp *http.Response)
	}{
		//{
		//	name: "right solution",
		//	solve: func(challenge string, difficulty int) (string, string) {
		//		var solution int
		//		for {
		//			solutionStr := fmt.Sprintf("%d", solution)
		//			hash := sha256.Sum256([]byte(challenge + solutionStr))
		//			hashStr := hex.EncodeToString(hash[:])
		//			if strings.HasPrefix(hashStr, strings.Repeat("0", difficulty)) {
		//				return challenge, solutionStr
		//			}
		//			solution++
		//		}
		//	},
		//	check: func(resp *http.Response) {
		//		require.Equal(t, http.StatusOK, resp.StatusCode)
		//		defer func() {
		//			_ = resp.Body.Close()
		//		}()
		//
		//		payload, err := io.ReadAll(resp.Body)
		//		require.NoError(t, err)
		//
		//		var realResp response.SolveResponse
		//		err = jsoniter.Unmarshal(payload, &realResp)
		//		require.NoError(t, err)
		//
		//		require.NotEmpty(t, realResp.Quote)
		//		require.True(t, slices.Contains(cfg.QuotesList, realResp.Quote))
		//	},
		//},
		//{
		//	name: "wrong solution",
		//	solve: func(challenge string, difficulty int) (string, string) {
		//		return challenge, "wrong"
		//	},
		//	check: func(resp *http.Response) {
		//		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		//		defer func() {
		//			_ = resp.Body.Close()
		//		}()
		//
		//		payload, err := io.ReadAll(resp.Body)
		//		require.NoError(t, err)
		//
		//		var realResp api.ErrorResponse
		//		err = jsoniter.Unmarshal(payload, &realResp)
		//		require.NoError(t, err)
		//
		//		require.Equal(t, "incorrect solution", realResp.Error.Text)
		//	},
		//},
		//{
		//	name: "wrong challenge",
		//	solve: func(challenge string, difficulty int) (string, string) {
		//		return "wrong", "wrong"
		//	},
		//	check: func(resp *http.Response) {
		//		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		//		defer func() {
		//			_ = resp.Body.Close()
		//		}()
		//
		//		payload, err := io.ReadAll(resp.Body)
		//		require.NoError(t, err)
		//
		//		var realResp api.ErrorResponse
		//		err = jsoniter.Unmarshal(payload, &realResp)
		//		require.NoError(t, err)
		//
		//		require.Equal(t, "challenge not found", realResp.Error.Text)
		//	},
		//},
		{
			name: "expired challenge",
			solve: func(challenge string, difficulty int) (string, string) {
				time.Sleep(cfg.Challenge.TTL)
				time.Sleep(time.Second)
				var solution int
				for {
					solutionStr := fmt.Sprintf("%d", solution)
					hash := sha256.Sum256([]byte(challenge + solutionStr))
					hashStr := hex.EncodeToString(hash[:])
					if strings.HasPrefix(hashStr, strings.Repeat("0", difficulty)) {
						return challenge, solutionStr
					}
					solution++
				}
			},
			check: func(resp *http.Response) {
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
				defer func() {
					_ = resp.Body.Close()
				}()

				payload, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var realResp api.ErrorResponse
				err = jsoniter.Unmarshal(payload, &realResp)
				require.NoError(t, err)

				require.Equal(t, "challenge not found", realResp.Error.Text)
			},
		},
	}

	url := fmt.Sprintf("%s/solve", baseURL)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			challenge, difficulty := getChallenge()
			challenge, solution := tc.solve(challenge, difficulty)

			payload := fmt.Sprintf("{\"challenge\": \"%s\", \"solution\": \"%s\"}", challenge, solution)
			resp, err := httpClient.Post(url, "application/json", strings.NewReader(payload))
			require.NoError(t, err)

			tc.check(resp)
		})
	}
}
