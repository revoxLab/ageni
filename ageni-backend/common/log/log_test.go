package log

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
	"unsafe"
)

func init() {
	cfg := &Config{
		AppName: "test",
		Debug:   false,
	}
	Init(cfg)
}

func TestDebug(t *testing.T) {
	Debug("hello debug")
	Debugf("hello number=%d", 100)
}

// go test -v -test.run TestInfo
func TestInfo(t *testing.T) {
	Info("hello")
	Infof("hello number=%d", 100)
}

func TestWarn(t *testing.T) {
	Warn("hello")
	Warnf("hello  number=%d", 100)
}

// go test -v -test.run TestError
func TestError(t *testing.T) {
	Error("hello")
	Errorf("hello number=%d", 100)
}

func TestFatal(t *testing.T) {
	Fatal("hello")
	Fatalf("hello  number=%d", 100)
}

func TestContextLogger(t *testing.T) {
	ctx := context.Background()
	fmt.Printf("%p, %v\n", logger(), unsafe.Sizeof(logger()))
	ctx = SetContextLog(ctx, logger())
	newLogger := GetFromContext(ctx)
	fmt.Printf("%p, %v\n", newLogger, unsafe.Sizeof(newLogger))
	newLogger.Info("test trace id")
	logger().Info("test trace id2")

	ginCtx := &gin.Context{}
	NewTraceIdWithGinContext(ginCtx)
	InfoWithGinContext(ginCtx, "InfoWithContext gin")

	ctx = NewTraceIdWithContext(ctx)
	InfoWithContext(ctx, "InfoWithContext golang")

	t.Log(ctx.Value(TraceIdKey))
}
