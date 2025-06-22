package observability

import (
	"context"
	"runtime/debug"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
)

var EnableGoroutineTracing = false

func Call(ctx context.Context, fn func(context.Context)) {
	defer func() { PanicIfNotNil(ctx, recover()) }()
	fn(ctx)
}

func CallSafe(ctx context.Context, fn func(context.Context)) {
	defer func() { ReportPanicIfNotNil(ctx, recover()) }()
	fn(ctx)
}

func Go(ctx context.Context, fn func(context.Context)) {
	if EnableGoroutineTracing {
		stack := debug.Stack()
		ctx = addGoroutinesStack(ctx, stack)
	}
	go Call(ctx, fn)
}

func GoSafe(ctx context.Context, fn func(context.Context)) {
	if EnableGoroutineTracing {
		stack := debug.Stack()
		ctx = addGoroutinesStack(ctx, stack)
	}
	go CallSafe(ctx, fn)
}

func GoSafeRestartable(ctx context.Context, fn func(context.Context)) {
	if EnableGoroutineTracing {
		stack := debug.Stack()
		ctx = addGoroutinesStack(ctx, stack)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			CallSafe(ctx, fn)
		}
	}()
}

func addGoroutinesStack(ctx context.Context, stack []byte) context.Context {
	belt.GetFields(ctx).ForEachField(func(f *field.Field) bool {
		if f.Key == "goroutines_stack" {
			prevStack, ok := f.Value.([]byte)
			if !ok {
				return true
			}
			stack = append(stack, prevStack...)
			return false
		}
		return true
	})
	return belt.WithField(ctx, "goroutines_stack", string(stack))
}
