package dbaccess

import "log/slog"

func LogDbResult(queryName string, result any, err error) {
	if err != nil {
		slog.Error(queryName, "error", err)
	} else {
		slog.Info(queryName, "result", result)
	}
}
