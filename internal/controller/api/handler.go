package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/AydarZaynutdinov/pow/internal/controller/api/utils"
	"github.com/AydarZaynutdinov/pow/internal/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	ContentTypeHeader = "Content-Type"
)

type (
	Error struct {
		Text string `json:"text"`
	}
	ErrorResponse struct {
		Error Error `json:"error"`
	}
)

type Handler struct {
	logger logger.Logger
}

func NewHandler(logger logger.Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) InternalServerError(ctx context.Context, w http.ResponseWriter) error {
	resp := ErrorResponse{
		Error: Error{
			Text: "internal server error",
		},
	}
	return h.WriteResponse(ctx, w, resp, http.StatusInternalServerError)
}

func (h *Handler) Handle(handler func(r *http.Request) (*utils.ControllerResult, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		ctx := r.Context()
		ctx = logger.AddRequestIDToContext(ctx, requestID)

		log := h.logger.WithContext(ctx)
		log.Info("started incoming http request")

		result, err := handler(r.WithContext(ctx))
		if err != nil {
			log.Error("failed incoming http request", zap.Error(err))

			if errors.Is(err, context.Canceled) {
				resp := ErrorResponse{
					Error: Error{
						Text: "canceled",
					},
				}
				_ = h.WriteResponse(ctx, w, resp, http.StatusInternalServerError)
			} else if errors.Is(err, context.DeadlineExceeded) {
				resp := ErrorResponse{
					Error: Error{
						Text: "deadline exceeded",
					},
				}
				_ = h.WriteResponse(ctx, w, resp, http.StatusInternalServerError)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				var typedErr *utils.ControllerStandardError
				switch {
				case errors.As(err, &typedErr):
					resp := ErrorResponse{
						Error: Error{
							Text: typedErr.Error(),
						},
					}
					_ = h.WriteResponse(ctx, w, resp, typedErr.HTTPCode)
				default:
					resp := ErrorResponse{
						Error: Error{
							Text: err.Error(),
						},
					}
					_ = h.WriteResponse(ctx, w, resp, http.StatusInternalServerError)
					log.Error("unknown handler error", zap.Error(err))
				}
			}
			return
		}

		err = h.WriteResponse(ctx, w, result.Response, result.HTTPCode)
		if err == nil {
			log.Info("completed incoming http request", zap.Int("status", result.HTTPCode))
		}
	}
}

func (h *Handler) WriteResponse(
	ctx context.Context,
	w http.ResponseWriter,
	response any,
	httpStatus int,
) error {
	encoder := utils.JSON
	w.Header().Set(ContentTypeHeader, encoder.ContentType())
	w.WriteHeader(httpStatus)
	if err := encoder.Encode(w, response); err != nil {
		h.logger.WithContext(ctx).Error("failed to write response")
		return err
	}

	return nil
}
