package dbaccess

import "log/slog"

func LogDbResult(q Query, result any, err error) {
	if err != nil {
		slog.Error(q.Name, "error", err)
	} else {
		slog.Info(q.Name, "result", result)
	}
}
