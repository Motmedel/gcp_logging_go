package gcp_logging

import (
	"context"
	"encoding/json"
	"fmt"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	motmedelHttpContext "github.com/Motmedel/utils_go/pkg/http/context"
	motmedelHttpTypes "github.com/Motmedel/utils_go/pkg/http/types"
	motmedelLog "github.com/Motmedel/utils_go/pkg/log"
	"io"
	"log/slog"
)

func ReplaceAttr(groups []string, attr slog.Attr) slog.Attr {
	if len(groups) > 0 {
		return attr
	}

	switch attr.Key {
	case slog.TimeKey:
		attr.Key = "time"
	case slog.LevelKey:
		attr.Key = "severity"
	case slog.MessageKey:
		attr.Key = "message"
	case slog.SourceKey:
		if source, ok := attr.Value.Any().(*slog.Source); ok {
			return slog.Group(
				"logging.googleapis.com/sourceLocation",
				"file", source.File,
				"line", source.Line,
				"function", source.Function,
			)
		}
	}

	return attr
}

func ExtractHttpContext(ctx context.Context, record *slog.Record) error {
	if record == nil {
		return nil
	}

	if httpContext, ok := ctx.Value(motmedelHttpContext.HttpContextContextKey).(*motmedelHttpTypes.HttpContext); ok {
		if logEntry := ParseHttp(httpContext.Request, httpContext.Response); logEntry != nil {
			logEntryBytes, err := json.Marshal(logEntry)
			if err != nil {
				return motmedelErrors.MakeErrorWithStackTrace(
					fmt.Errorf("json marshal (http context log entry): %w", err),
					logEntry,
				)
			}

			var logEntryMap map[string]any
			if err = json.Unmarshal(logEntryBytes, &logEntryMap); err != nil {
				return motmedelErrors.MakeErrorWithStackTrace(
					fmt.Errorf("json unmarshal (http context log entry map): %w", err),
					logEntry,
				)
			}

			record.Add(motmedelLog.AttrsFromMap(logEntryMap)...)
		}
	}

	return nil
}

var HttpContextExtractor = motmedelLog.ContextExtractorFunction(ExtractHttpContext)

func MakeLogger(level slog.Leveler, writer io.Writer) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			writer,
			&slog.HandlerOptions{
				AddSource:   true,
				Level:       level,
				ReplaceAttr: ReplaceAttr,
			},
		),
	)
}
