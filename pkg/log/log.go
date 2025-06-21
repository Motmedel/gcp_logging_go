package log

import (
	"context"
	"fmt"
	"github.com/Motmedel/gcp_logging_go/pkg/gcp_logging"
	gcpLoggingTypes "github.com/Motmedel/gcp_logging_go/pkg/types"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	motmedelHttpContext "github.com/Motmedel/utils_go/pkg/http/context"
	motmedelHttpTypes "github.com/Motmedel/utils_go/pkg/http/types"
	motmedelJson "github.com/Motmedel/utils_go/pkg/json"
	motmedelLog "github.com/Motmedel/utils_go/pkg/log"
	motmedelLogHandler "github.com/Motmedel/utils_go/pkg/log/handler"
	"io"
	"log/slog"
)

func ParseHttpContext(httpContext *motmedelHttpTypes.HttpContext) *gcpLoggingTypes.LogEntry {
	if httpContext == nil {
		return nil
	}

	return gcp_logging.ParseHttp(httpContext.Request, httpContext.Response)
}

type HttpContextExtractor struct {
}

func (httpContextExtractor *HttpContextExtractor) Handle(ctx context.Context, record *slog.Record) error {
	if record == nil {
		return nil
	}

	if httpContext, ok := ctx.Value(motmedelHttpContext.HttpContextContextKey).(*motmedelHttpTypes.HttpContext); ok && httpContext != nil {
		if logEntry := ParseHttpContext(httpContext); logEntry != nil {
			logEntryMap, err := motmedelJson.ObjectToMap(logEntry)
			if err != nil {
				return motmedelErrors.New(fmt.Errorf("object to map: %w", err), logEntry)
			}

			record.Add(motmedelLog.AttrsFromMap(logEntryMap)...)
		}
	}

	return nil
}

func LoggerReplaceAttr(groups []string, attr slog.Attr) slog.Attr {
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

func MakeLogger(level slog.Leveler, writer io.Writer) *slog.Logger {
	return slog.New(
		motmedelLogHandler.New(
			slog.NewJSONHandler(
				writer,
				&slog.HandlerOptions{AddSource: true, Level: level, ReplaceAttr: LoggerReplaceAttr},
			),
		),
	)
}
