package logger

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Key string

const (
	ContextLogger Key = "logger-value"
	SessionId     Key = "session-id"
	BotChat       Key = "bot-chat"
	BotMessage    Key = "bot-message"
)

func Context(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ContextLogger, l)
}

func AddId(ctx context.Context, chat int64, msg int) (context.Context, string) {
	v := uuid.New()
	return context.WithValue(context.WithValue(context.WithValue(ctx, SessionId, v.String()), BotChat, chat), BotMessage, msg), v.String()
}

func FromCtx(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ContextLogger).(*zap.Logger)
	if !ok {
		return nil
	}
	return logger
}

func SessionIdFromCtx(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(SessionId).(string)
	return v, ok
}

func BotChatFromCtx(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(BotChat).(int64)
	return v, ok
}

func BotMessageFromCtx(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(BotMessage).(int)
	return v, ok
}

func Log(ctx context.Context, lvl zapcore.Level, msg string, fields ...zap.Field) bool {
	logger := FromCtx(ctx)
	if logger == nil {
		return false
	}

	if session, ok := SessionIdFromCtx(ctx); ok {
		fields = append(fields, zap.String(string(SessionId), session))
	}
	if botChat, ok := BotChatFromCtx(ctx); ok {
		fields = append(fields, zap.Int64(string(BotChat), botChat))
	}
	if botMessage, ok := BotMessageFromCtx(ctx); ok {
		fields = append(fields, zap.Int(string(BotMessage), botMessage))
	}

	logger.Log(lvl, msg, fields...)
	return true
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.DebugLevel, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.InfoLevel, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.WarnLevel, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.ErrorLevel, msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.FatalLevel, msg, fields...)
}

func DPanic(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.DPanicLevel, msg, fields...)
}

func Panic(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.PanicLevel, msg, fields...)
}
