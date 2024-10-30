package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type ChallengeResponse struct {
	Challenge  string `json:"challenge"`
	Difficulty int    `json:"difficulty"`
}

type SolveResponse struct {
	Quote string `json:"quote"`
}

type SolveRequest struct {
	Challenge string `json:"challenge"`
	Solution  string `json:"solution"`
}

var serverURL string

func init() {
	// parsing config file path
	flagServer := flag.String("server", "http://127.0.0.1:8080", "server base URL")
	flag.Parse()

	serverURL = *flagServer
}

func main() {
	for {
		time.Sleep(30 * time.Second)
		challengeResp, err := getChallenge()
		if err != nil {
			fmt.Println(err)
			continue
		}

		solution := solve(challengeResp.Challenge, challengeResp.Difficulty)

		quote, err := checkSolution(challengeResp.Challenge, solution)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("quote is: %s\n", quote)
	}
}

func getChallenge() (*ChallengeResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/pow/challenge", serverURL))
	if err != nil {
		return nil, fmt.Errorf("challenge: error during requesting challenge: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("challenge: status is not 200: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("challenge: error reading challenge response: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var challengeResp ChallengeResponse
	err = jsoniter.Unmarshal(body, &challengeResp)
	if err != nil {
		return nil, fmt.Errorf("challenge: error unmarshalling challenge response: %v", err)
	}

	return &challengeResp, nil
}

func solve(challenge string, difficulty int) string {
	var solution int
	for {
		solutionStr := fmt.Sprintf("%d", solution)
		hash := sha256.Sum256([]byte(challenge + solutionStr))
		hashStr := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hashStr, strings.Repeat("0", difficulty)) {
			return solutionStr
		}
		solution++
	}
}

func checkSolution(challenge, solution string) (string, error) {
	req := SolveRequest{
		Challenge: challenge,
		Solution:  solution,
	}
	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("solve: error marshalling solve request: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/pow/solve", serverURL), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("solve: error sending solve request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("solve: response status is not 200: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("solve: error reading challenge response: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var solveResponse SolveResponse
	err = jsoniter.Unmarshal(body, &solveResponse)
	if err != nil {
		return "", fmt.Errorf("solve: error unmarshalling solve response: %v", err)
	}

	return solveResponse.Quote, nil
}
