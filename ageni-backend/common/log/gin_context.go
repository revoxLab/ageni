package log

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func SetTraceIdWithGinContext(ctx *gin.Context, traceId string) {
	ctx.Set(TraceIdKey, traceId)
}

func NewTraceIdWithGinContext(ctx *gin.Context) {
	ctx.Set(TraceIdKey, uuid.NewV4().String())
}

func InfoWithGinContext(ctx *gin.Context, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceIdKey)).Info(args...)
}

func DebugfWithGinContext(ctx *gin.Context, msg string, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceIdKey)).Debugf(msg, args...)
}

func DebugWithGinContext(ctx *gin.Context, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceIdKey)).Debug(args...)
}

func InfofWithGinContext(ctx *gin.Context, msg string, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceIdKey)).Infof(msg, args...)
}

func ErrorfWithGinContext(ctx *gin.Context, msg string, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceIdKey)).Errorf(msg, args...)
}

func ErrorWithGinContext(ctx *gin.Context, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceIdKey)).Error(args...)
}

func WarnfWithGinContext(ctx *gin.Context, msg string, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceIdKey)).Warnf(msg, args...)
}

func WarnWithGinContext(ctx *gin.Context, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceIdKey)).Warn(args...)
}

func FatalfWithGinContext(ctx *gin.Context, msg string, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceIdKey)).Fatalf(msg, args...)
}

func FatalWithGinContext(ctx *gin.Context, args ...interface{}) {
	logger().With(TraceIdKey, ctx.Value(TraceIdKey)).Fatal(args...)
}
