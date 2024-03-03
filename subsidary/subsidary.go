package subsidary

import (
	"context"
	"github.com/RustGrub/FunnyGoService/consts"
	"github.com/RustGrub/FunnyGoService/logger"
	"github.com/RustGrub/FunnyGoService/logger/std"
)

// NewLogUsingContext Просто побаловался
func NewLogUsingContext(ctx context.Context, msg, lvl string) {
	reqID := ctx.Value(consts.ReqID).(string)
	if reqID == "" {
		reqID = "Unknown Request"
	}

	l, ok := ctx.Value(consts.Logger).(logger.Logger)
	if !ok {
		return
	}
	switch lvl {
	case std.DebugLevel:
		l.Debug(reqID, msg)
	case std.FatalLevel:
		l.Fatal(reqID, msg)
	case std.ErrLevel:
		l.Error(reqID, msg)
	case std.InfoLevel:
		l.Info(reqID, msg)
	case std.WarnignLevel:
		l.Warning(reqID, msg)
	default:
		l.Info(reqID, msg)
	}
}
func NewContextWithLogger(ctx context.Context, l logger.Logger) context.Context {
	return context.WithValue(ctx, consts.Logger, l)
}
