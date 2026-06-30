package messaging

import (
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	"jolly/backend/common/log"
)

// LoggingMiddleware logs incoming messages, their duration, and any errors.
// It also extracts the correlation ID and sets it in the message context for downstream handlers.
func LoggingMiddleware(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		corrID := middleware.MessageCorrelationID(msg)
		if corrID == "" {
			corrID = log.CorrelationIDFromContext(msg.Context())
		}

		logger := slog.With(
			"correlation_id", corrID,
			"message_uuid", msg.UUID,
		)

		ctx := log.ToContext(msg.Context(), logger)
		ctx = log.ContextWithCorrelationID(ctx, corrID)
		msg.SetContext(ctx)

		logger.Info("Processing broker message")

		start := time.Now()
		events, err := h(msg)
		duration := time.Since(start)

		if err != nil {
			logger.Error("Broker message processing failed", "error", err, "duration", duration.String())
		} else {
			logger.Info("Broker message processed successfully", "duration", duration.String())
		}

		return events, err
	}
}
