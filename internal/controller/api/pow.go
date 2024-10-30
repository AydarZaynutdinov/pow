package api

import (
	"net/http"

	"github.com/AydarZaynutdinov/pow/internal/controller/api/request"
	"github.com/AydarZaynutdinov/pow/internal/controller/api/response"
	"github.com/AydarZaynutdinov/pow/internal/controller/api/utils"
)

func (c *Controller) GetChallenge(r *http.Request) (*utils.ControllerResult, error) {
	ctx := r.Context()

	generatedChallenge, err := c.service.GenerateChallenge(ctx)
	if err != nil {
		return nil, utils.NewControllerStandardErrorInternalServerError(err)
	}

	return utils.NewControllerResultOK(
		response.ChallengeResponse{
			Challenge:  generatedChallenge,
			Difficulty: c.service.GetDifficulty(),
		},
	), nil
}

func (c *Controller) Solve(r *http.Request) (*utils.ControllerResult, error) {
	ctx := r.Context()

	solveChallengeRequest, err := request.NewSolveRequest(ctx, r, c.logger)
	if err != nil {
		return nil, utils.NewControllerStandardErrorBadRequest(err)
	}

	quote, err := c.service.CheckSolution(ctx, solveChallengeRequest.Challenge, solveChallengeRequest.Solution)
	if err != nil {
		if c.service.IsSolutionIncorrect(err) {
			return nil, utils.NewControllerStandardErrorBadRequest(err)
		}
		return nil, utils.NewControllerStandardErrorInternalServerError(err)
	}

	return utils.NewControllerResultOK(
		response.SolveResponse{
			Quote: quote,
		},
	), nil
}
