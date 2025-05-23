package xlogger

import (
	"context"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/facebookincubator/go-belt/tool/logger/adapter"
)

type EmitterWithMaxLoggingLevel struct {
	logger.Emitter
	MaxLevel logger.Level
}

var _ logger.Emitter = (*EmitterWithMaxLoggingLevel)(nil)

func WithMaxLoggingLevel(l logger.Logger, maxLevel logger.Level) logger.Logger {
	emitter := &EmitterWithMaxLoggingLevel{Emitter: l.Emitter(), MaxLevel: maxLevel}
	return adapter.LoggerFromEmitter(emitter)
}

func (e *EmitterWithMaxLoggingLevel) Emit(
	entry *logger.Entry,
) {
	switch entry.Level {
	case logger.LevelPanic, logger.LevelFatal:
		e.Emitter.Emit(entry)
		return
	}

	if entry.Level < e.MaxLevel {
		entry.Level = e.MaxLevel
	}
	e.Emitter.Emit(entry)
}

func CtxWithMaxLoggingLevel(ctx context.Context, maxLevel logger.Level) context.Context {
	emitter := &EmitterWithMaxLoggingLevel{Emitter: logger.FromCtx(ctx).Emitter(), MaxLevel: maxLevel}
	return logger.CtxWithLogger(ctx, adapter.LoggerFromEmitter(emitter))
}
