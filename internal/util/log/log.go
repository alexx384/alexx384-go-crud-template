package log

import (
	"context"
	"github.com/gin-gonic/gin"
	slogctx "github.com/veqryn/slog-context"
	"log/slog"
	"runtime"
	"time"
)

func Info(ctx *gin.Context, msg string, args ...any) {
	log(ctx.Request.Context(), slogctx.FromCtx(ctx), slog.LevelInfo, msg, 3, args...)
}

func Warn(ctx *gin.Context, msg string, args ...any) {
	log(ctx.Request.Context(), slogctx.FromCtx(ctx), slog.LevelWarn, msg, 3, args...)
}

func Error(ctx *gin.Context, msg string, args ...any) {
	log(ctx.Request.Context(), slogctx.FromCtx(ctx), slog.LevelError, msg, 3, args...)
}

func Debug(ctx *gin.Context, msg string, args ...any) {
	log(ctx.Request.Context(), slogctx.FromCtx(ctx), slog.LevelDebug, msg, 3, args...)
}

// log is the low-level logging method for methods that take ...any.
// It must always be called directly by an exported logging method
// or function, because it uses a fixed call depth to obtain the pc.
// This is copied from golang sdk.
func log(ctx context.Context, l *slog.Logger, level slog.Level, msg string, skipCalls int, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}

	if !l.Enabled(ctx, level) {
		return
	}

	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(skipCalls, pcs[:])
	pc = pcs[0]

	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	_ = l.Handler().Handle(ctx, r)
}
