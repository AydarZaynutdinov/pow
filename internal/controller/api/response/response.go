package response

type ChallengeResponse struct {
	Challenge  string `json:"challenge"`
	Difficulty int    `json:"difficulty"`
}

type SolveResponse struct {
	Quote string `json:"quote"`
}
