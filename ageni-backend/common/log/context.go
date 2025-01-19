package log

import (
	"context"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

const TraceIdKey = "__trace_id__"

type TraceLogKey string

func SetTraceIdWithContext(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, TraceLogKey(TraceIdKey), traceId)
}

func NewTraceIdWithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, TraceLogKey(TraceIdKey), uuid.NewV4().String())
}

func InfoWithContext(ctx context.Context, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceLogKey(TraceIdKey))).Info(args...)
}

func DebugfWithContext(ctx context.Context, msg string, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceLogKey(TraceIdKey))).Debugf(msg, args...)
}

func DebugWithContext(ctx context.Context, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceLogKey(TraceIdKey))).Debug(args...)
}

func InfofWithContext(ctx context.Context, msg string, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceLogKey(TraceIdKey))).Infof(msg, args...)
}

func ErrorfWithContext(ctx context.Context, msg string, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceLogKey(TraceIdKey))).Errorf(msg, args...)
}

func ErrorWithContext(ctx context.Context, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceLogKey(TraceIdKey))).Error(args...)
}

func WarnfWithContext(ctx context.Context, msg string, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceLogKey(TraceIdKey))).Warnf(msg, args...)
}

func WarnWithContext(ctx context.Context, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceLogKey(TraceIdKey))).Warn(args...)
}

func FatalfWithContext(ctx context.Context, msg string, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceLogKey(TraceIdKey))).Fatalf(msg, args...)
}

func FatalWithContext(ctx context.Context, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceLogKey(TraceIdKey))).Fatal(args...)
}

func SetContextLog(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	setup := logger
	setup = setup.With(TraceIdKey, uuid.NewV4().String())
	//fmt.Printf("%p, %v\n", setup, unsafe.Sizeof(setup))
	return context.WithValue(ctx, TraceLogKey(TraceIdKey), setup)
}

func GetFromContext(ctx context.Context) *zap.SugaredLogger {
	if l, ok := ctx.Value(TraceLogKey(TraceIdKey)).(*zap.SugaredLogger); ok {
		return l
	}
	return zap.S()
}
