package observability

import (
	"context"
	"runtime/debug"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
)

func Call(ctx context.Context, fn func(context.Context)) {
	defer func() { PanicIfNotNil(ctx, recover()) }()
	fn(ctx)
}

func CallSafe(ctx context.Context, fn func(context.Context)) {
	defer func() { ReportPanicIfNotNil(ctx, recover()) }()
	fn(ctx)
}

func Go(ctx context.Context, fn func(context.Context)) {
	stack := debug.Stack()
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
	ctx = belt.WithField(ctx, "goroutines_stack", string(stack))
	go Call(ctx, fn)
}

func GoSafe(ctx context.Context, fn func(context.Context)) {
	stack := debug.Stack()
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
	ctx = belt.WithField(ctx, "goroutines_stack", string(stack))
	go CallSafe(ctx, fn)
}
