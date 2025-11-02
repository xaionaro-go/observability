package observability

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/tool/experimental/errmon"
	"github.com/facebookincubator/go-belt/tool/logger"
)

func PanicIfNotNil(ctx context.Context, r any) {
	if r == nil {
		return
	}
	Panic(ctx, r)
}

func Panic(ctx context.Context, r any) {
	ReportPanic(ctx, r)
	time.Sleep(time.Second)
	panic(fmt.Sprintf("%#+v", r))
}

func ReportPanicIfNotNil(ctx context.Context, r any) bool {
	if r == nil {
		return false
	}
	return ReportPanic(ctx, r)
}

func ReportPanic(ctx context.Context, r any) bool {
	logger.FromCtx(ctx).
		WithField("error_event_exception_stack_trace", string(debug.Stack())).
		Errorf("got panic: %v", r)
	errmon.ObserveRecoverCtx(ctx, r)
	belt.Flush(ctx)
	return true
}
