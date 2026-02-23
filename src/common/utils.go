package common

import (
	"FGG-Service/api/generated/auth"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SendJSONErrorResponse(ctx echo.Context, err error) error {
	var badRequestError *BadRequestError
	var unauthorizedError *UnauthorizedError
	var notFoundError *NotFoundError
	var conflictError *ConflictError
	var unprocessableError *UnprocessableError

	apiCode := http.StatusInternalServerError

	switch {
	case errors.As(err, &badRequestError):
		apiCode = http.StatusBadRequest
	case errors.As(err, &unauthorizedError):
		apiCode = http.StatusUnauthorized
	case errors.As(err, &notFoundError):
		apiCode = http.StatusNotFound
	case errors.As(err, &conflictError):
		apiCode = http.StatusConflict
	case errors.As(err, &unprocessableError):
		apiCode = http.StatusUnprocessableEntity
	}

	apiError := convertToError(err)

	return ctx.JSON(apiCode, apiError)
}

func convertToError(err error) auth.Error {
	var appError AppError
	if errors.As(err, &appError) {
		return auth.Error{
			Code:    appError.GetCode(),
			Message: appError.GetMessage(),
		}
	}

	return auth.Error{
		Code:    "UNEXPECTED",
		Message: err.Error(),
	}
}

func DurationToISO8601(duration time.Duration) string {
	if duration == 0 {
		return "PT0S"
	}

	sign := ""
	if duration < 0 {
		sign = "-"
		duration = -duration
	}

	totalSeconds := int64(duration.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	result := fmt.Sprintf("%sPT", sign)

	if hours > 0 {
		result += fmt.Sprintf("%dH", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dM", minutes)
	}
	if seconds > 0 {
		result += fmt.Sprintf("%dS", seconds)
	}

	return result

}
