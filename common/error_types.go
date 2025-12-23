package common

import (
	"fmt"
)

type AppError interface {
	error
	GetCode() string
	GetMessage() string
}

type BaseError struct {
	Code    string
	Message string
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *BaseError) GetCode() string {
	return e.Code
}

func (e *BaseError) GetMessage() string {
	return e.Message
}

type AuthError struct {
	BaseError
}

type CookieNotFoundAuthError struct {
	AuthError
}

func NewCookieNotFoundAuthError() *CookieNotFoundAuthError {
	return &CookieNotFoundAuthError{
		AuthError{
			BaseError{
				Code: "COOKIE_NOT_FOUND",
				// TODO: Make a message
				Message: "",
			},
		},
	}
}

type ActiveSessionNotFoundAuthError struct {
	AuthError
}

func NewActiveSessionNotFoundAuthError() *ActiveSessionNotFoundAuthError {
	return &ActiveSessionNotFoundAuthError{
		AuthError{
			BaseError{
				Code: "ACTIVE_SESSION_NOT_FOUND",
				// TODO: Make a message
				Message: "",
			},
		},
	}
}

type WrongDataAuthError struct {
	AuthError
}

func NewWrongDataAuthError() *WrongDataAuthError {
	return &WrongDataAuthError{
		AuthError{
			BaseError{
				Code: "WRONG_AUTH_DATA",
				// TODO: Make a message
				Message: "",
			},
		},
	}
}

type SessionAlreadyExistsAuthError struct {
	AuthError
}

func NewSessionAlreadyExistsAuthError() *SessionAlreadyExistsAuthError {
	return &SessionAlreadyExistsAuthError{
		AuthError{
			BaseError{
				Code: "SESSION_ALREADY_EXISTS",
				// TODO: Make a message
				Message: "",
			},
		},
	}
}

type DatabaseError struct {
	BaseError
}

type NotFoundError struct {
	BaseError
}

type CurrentGameNotFoundError struct {
	NotFoundError
}

func NewCurrentGameNotFoundError() *CurrentGameNotFoundError {
	return &CurrentGameNotFoundError{
		NotFoundError{
			BaseError{
				Code:    "CURRENT_GAME_NOT_FOUND",
				Message: "The user doesn't have a current game. Roll the game to get one.",
			},
		},
	}
}

type CompletedTimersNotFoundError struct {
	NotFoundError
}

func NewCompletedTimersNotFoundError() *CompletedTimersNotFoundError {
	return &CompletedTimersNotFoundError{
		NotFoundError{
			BaseError{
				Code:    "COMPLETED_TIMERS_NOT_FOUND",
				Message: "The user doesn't have completed timers. Complete at least one timer to finish the game.",
			},
		},
	}
}

type UnplayedGamesNotFoundError struct {
	NotFoundError
}

func NewUnplayedGamesNotFoundError() *UnplayedGamesNotFoundError {
	message := fmt.Sprintf(
		"The user doesn't have unplayed games. Add at least %d to roll the game.",
		MinimumNumberOfUnplayedGames)

	return &UnplayedGamesNotFoundError{
		NotFoundError{
			BaseError{
				Code:    "UNPLAYED_GAMES_NOT_FOUND",
				Message: message,
			},
		},
	}
}

type CurrentTimerNotFoundError struct {
	NotFoundError
}

func NewCurrentTimerNotFoundError() *CurrentTimerNotFoundError {
	return &CurrentTimerNotFoundError{
		NotFoundError{
			BaseError{
				Code:    "CURRENT_TIMER_NOT_FOUND",
				Message: "The user doesn't have a current timer. Create a timer so you can control it.",
			},
		},
	}
}

type AvailableRollsNotFoundError struct {
	NotFoundError
}

func NewAvailableRollsNotFoundError() *AvailableRollsNotFoundError {
	return &AvailableRollsNotFoundError{
		NotFoundError{
			BaseError{
				Code:    "AVAILABLE_ROLLS_NOT_FOUND",
				Message: "The user doesn't have available rolls. Complete the timer to get one.",
			},
		},
	}
}

type ConflictStateError struct {
	BaseError
}

type CurrentTimerIncorrectStateError struct {
	ConflictStateError
}

func NewCurrentTimerIncorrectStateError(timerState TimerStateType) *CurrentTimerIncorrectStateError {
	message := fmt.Sprintf(
		"This action cannot be performed. The current timer is in the \"%s\" state.",
		timerState)

	return &CurrentTimerIncorrectStateError{
		ConflictStateError{
			BaseError{
				Code:    "CURRENT_TIMER_INCORRECT_STATE",
				Message: message,
			},
		},
	}
}

type CurrentGameAlreadyExistsError struct {
	ConflictStateError
}

func NewCurrentGameAlreadyExistsError() *CurrentGameAlreadyExistsError {
	return &CurrentGameAlreadyExistsError{
		ConflictStateError{
			BaseError{
				Code:    "CURRENT_GAME_ALREADY_EXISTS",
				Message: "The current game has already been rolled. Complete it before you can roll a new one.",
			},
		},
	}
}

type UserNameAlreadyExistsError struct {
	ConflictStateError
}

func NewUserNameAlreadyExistsError() *UserNameAlreadyExistsError {
	return &UserNameAlreadyExistsError{
		ConflictStateError{
			BaseError{
				Code: "USER_NAME_ALREADY_EXISTS",
				// TODO: Make a message
				Message: "",
			},
		},
	}
}

type UserEmailAlreadyExistsError struct {
	ConflictStateError
}

func NewUserEmailAlreadyExistsError() *UserEmailAlreadyExistsError {
	return &UserEmailAlreadyExistsError{
		ConflictStateError{
			BaseError{
				Code: "USER_EMAIL_ALREADY_EXISTS",
				// TODO: Make a message
				Message: "",
			},
		},
	}
}
