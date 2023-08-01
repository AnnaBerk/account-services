package sl

import (
	"fmt"
	"golang.org/x/exp/slog"
)

func Err(err error) slog.Attr {
	e, ok := err.(fmt.Stringer)
	if !ok {
		return slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		}
	}
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(e.String()),
	}
}
