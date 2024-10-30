package logger

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

const (
	unknownPath = "unknown path"

	logFieldKeyRequestID   = "request_id"
	logFieldKeyRequestPath = "request_path"
)

type contextLogFieldsType struct{}

var contextLogFieldsKey contextLogFieldsType

func AddRequestIDToContext(ctx context.Context, uuid string) context.Context {
	return addLogFieldsToContext(ctx, zap.String(logFieldKeyRequestID, uuid))
}

func AddRequestInfoToContext(ctx context.Context, request *http.Request) context.Context {
	path := unknownPath
	if request.URL != nil {
		path = request.URL.Path
	}

	return addLogFieldsToContext(ctx, zap.String(logFieldKeyRequestPath, path))
}

func addLogFieldsToContext(ctx context.Context, newFields ...zap.Field) context.Context {
	fields, hasFields := ctx.Value(contextLogFieldsKey).([]zap.Field)
	if !hasFields {
		if len(newFields) == 0 {
			return context.Background()
		}

		return context.WithValue(ctx, contextLogFieldsKey, newFields)
	}

	return context.WithValue(ctx, contextLogFieldsKey, append(fields, newFields...))
}

func getLogFieldsFromContext(ctx context.Context) []zap.Field {
	fields, hasFields := ctx.Value(contextLogFieldsKey).([]zap.Field)
	if !hasFields {
		return nil
	}

	return fields
}
