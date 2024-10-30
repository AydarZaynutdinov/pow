package request

import (
	"context"
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/AydarZaynutdinov/pow/internal/logger"
)

type SolveRequest struct {
	Challenge string `json:"challenge"`
	Solution  string `json:"solution"`
}

func NewSolveRequest(
	ctx context.Context,
	request *http.Request,
	logger logger.Logger,
) (*SolveRequest, error) {
	bodyPayload, err := io.ReadAll(request.Body)
	if err != nil {
		logger.WithContext(ctx).Error("failed to read request body", zap.Error(err))
		return nil, fmt.Errorf("failed to read request body")
	}
	_ = request.Body.Close()

	eventRequest := &SolveRequest{}
	err = jsoniter.Unmarshal(bodyPayload, eventRequest)
	if err != nil {
		logger.WithContext(ctx).Error("failed to unmarshal request body", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal request body")
	}

	return eventRequest, nil
}
