package service

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
)

const (
	ChallengeLen = 16
	SolutionLen  = 8
)

type PoW struct {
	random     *rand.Rand
	complexity uint8
}

func NewPoW(complexity uint8) *PoW {
	return &PoW{
		random:     rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
		complexity: complexity,
	}
}

func (pow *PoW) GenerateChallenge() []byte {
	challenge := make([]byte, ChallengeLen)
	_, _ = pow.random.Read(challenge)
	return challenge
}

func (pow *PoW) VerifySolution(challenge []byte, solution []byte) bool {
	hash := sha256.Sum256(append(challenge, solution...))
	hashHex := hex.EncodeToString(hash[:])
	requiredPrefix := strings.Repeat("0", int(pow.complexity))
	return strings.HasPrefix(hashHex, requiredPrefix)
}

func (pow *PoW) Solve(challenge []byte) []byte {
	solution := make([]byte, SolutionLen)
	for {
		_, err := pow.random.Read(solution)
		if err != nil {
			continue
		}

		if pow.VerifySolution(challenge, solution) {
			break
		}
	}
	return solution
}
