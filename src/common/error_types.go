package common

import (
	"FGG-Service/src/timers/types"
	"fmt"
	"strings"
)

type AppError interface {
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

type BadRequestError struct {
	*BaseError
}

func NewBadRequestError(message string) error {
	return &BadRequestError{
		&BaseError{
			Code:    "BAD_REQUEST",
			Message: message,
		},
	}
}

type UnauthorizedError struct {
	*BaseError
}

func NewCookieNotFoundUnauthorizedError() error {
	return &UnauthorizedError{
		&BaseError{
			Code:    "COOKIE_NOT_FOUND",
			Message: "Unable to retrieve the authentication cookie. Try logging in.",
		},
	}
}

func NewActiveSessionNotFoundUnauthorizedError() error {
	return &UnauthorizedError{
		&BaseError{
			Code:    "ACTIVE_SESSION_NOT_FOUND",
			Message: "An authentication session couldn't be found. Try logging in.",
		},
	}
}

type NotFoundError struct {
	*BaseError
}

func NewCurrentGameNotFoundError() error {
	return &NotFoundError{
		&BaseError{
			Code:    "CURRENT_GAME_NOT_FOUND",
			Message: "The user doesn't have a current game. Roll the game to get one.",
		},
	}
}

func NewCompletedTimersNotFoundError() error {
	return &NotFoundError{
		&BaseError{
			Code:    "COMPLETED_TIMERS_NOT_FOUND",
			Message: "The user doesn't have completed timers. Complete at least one timer to finish the game.",
		},
	}
}

func NewUnplayedGamesNotFoundError() error {
	message := fmt.Sprintf(
		"The user doesn't have unplayed games. Add at least %d to roll the game.",
		MinimumNumberOfUnplayedGames)

	return &NotFoundError{
		&BaseError{
			Code:    "UNPLAYED_GAMES_NOT_FOUND",
			Message: message,
		},
	}
}

func NewCurrentTimerNotFoundError() error {
	return &NotFoundError{
		&BaseError{
			Code:    "CURRENT_TIMER_NOT_FOUND",
			Message: "The user doesn't have a current timer. Create a timer so you can control it.",
		},
	}
}

func NewAvailableRollsNotFoundError() error {
	return &NotFoundError{
		&BaseError{
			Code:    "AVAILABLE_ROLLS_NOT_FOUND",
			Message: "The user doesn't have available rolls. Complete the timer to get one.",
		},
	}
}

func NewUserLoginNotFoundError(userLogin string) error {
	message := fmt.Sprintf(
		"The user login \"%s\" wasn't found.",
		userLogin)

	return &NotFoundError{
		&BaseError{
			Code:    "USER_LOGIN_NOT_FOUND",
			Message: message,
		},
	}
}

func NewLastWheelEffectsNotFoundError() error {
	return &NotFoundError{
		&BaseError{
			Code:    "LAST_WHEEL_EFFECTS_NOT_FOUND",
			Message: "The last wheel effects not found. Try to roll wheel effects.",
		},
	}
}

func NewDisplayNameNotFoundError() error {
	return &NotFoundError{
		&BaseError{
			Code:    "DISPLAY_NAME_NOT_FOUND",
			Message: "The user have not added a display name.",
		},
	}
}

func NewWheelEffectNameNotFoundError() error {
	return &NotFoundError{
		&BaseError{
			Code:    "WHEEL_EFFECT_NAME_NOT_FOUND",
			Message: "The wheel effect name not found. Change it or try again later.",
		},
	}
}

type ConflictError struct {
	*BaseError
}

func NewSessionAlreadyExistsConflictError() error {
	return &ConflictError{
		&BaseError{
			Code:    "SESSION_ALREADY_EXISTS",
			Message: "You're already logged in.",
		},
	}
}

func NewCurrentTimerIncorrectStateConflictError(timerState typetimers.TimerStateType) error {
	message := fmt.Sprintf(
		"This action cannot be performed. The current timer is in the \"%s\" state.",
		timerState)

	return &ConflictError{
		&BaseError{
			Code:    "CURRENT_TIMER_INCORRECT_STATE",
			Message: message,
		},
	}
}

func NewWishlistGameAlreadyExistsConflictError(gameName string) error {
	message := fmt.Sprintf(
		"The unplayed game \"%s\" has already been added.",
		gameName)

	return &ConflictError{
		&BaseError{
			Code:    "UNPLAYED_GAME_ALREADY_EXISTS",
			Message: message,
		},
	}
}

func NewCurrentGameAlreadyExistsConflictError() error {
	return &ConflictError{
		&BaseError{
			Code:    "CURRENT_GAME_ALREADY_EXISTS",
			Message: "The current game has already been rolled. Complete it before you can roll a new one.",
		},
	}
}

func NewUserNameAlreadyExistsConflictError() error {
	return &ConflictError{
		&BaseError{
			Code:    "USER_NAME_ALREADY_EXISTS",
			Message: "This username is already taken. Try another one.",
		},
	}
}

func NewUserEmailAlreadyExistsConflictError() error {
	return &ConflictError{
		&BaseError{
			Code:    "USER_EMAIL_ALREADY_EXISTS",
			Message: "This email is already taken. Try another one.",
		},
	}
}

func NewAvailableRollsExistConflictError() error {
	return &ConflictError{
		&BaseError{
			Code:    "AVAILABLE_ROLLS_EXIST",
			Message: "You have available rolls. You need to use them.",
		},
	}
}

func NewNotEnoughAvailableWheelEffectsConflictError() error {
	return &ConflictError{
		&BaseError{
			Code:    "NOT_ENOUGH_AVAILABLE_WHEEL_EFFECTS",
			Message: "You have not enough available wheel effects to roll. Contact with the administrator.",
		},
	}
}

func NewWrongDesiredChangeValueConflictError(changeSource string, changeSourceRule string) error {
	message := fmt.Sprintf(
		"The desired change value is invalid. The value for '%s' should be %s.",
		changeSource,
		changeSourceRule)

	return &ConflictError{
		&BaseError{
			Code:    "WRONG_DESIRED_CHANGE_VALUE",
			Message: message,
		},
	}
}

func NewNotEnoughCurrentPointsConflictError(changeSource string, requiredPoints int) error {
	message := fmt.Sprintf(
		"Not enough points for '%s' change source. It should be at least %d.",
		changeSource,
		requiredPoints)

	return &ConflictError{
		&BaseError{
			Code:    "NOT_ENOUGH_CURRENT_POINTS",
			Message: message,
		},
	}
}

type UnprocessableError struct {
	*BaseError
}

func NewWrongDataUnprocessableError() error {
	return &UnprocessableError{
		&BaseError{
			Code:    "WRONG_AUTH_DATA",
			Message: "Incorrect login or password. Try again.",
		},
	}
}

func NewUserNameUnprocessableError(name string, messageDetails string) error {
	message := fmt.Sprintf(
		"'%s' does not match the format. %s",
		name,
		messageDetails)

	return &UnprocessableError{
		&BaseError{
			Code:    "INCORRECT_USER_NAME_FORMAT",
			Message: message,
		},
	}
}

func NewNameUnprocessableError(name string, messageDetails string) error {
	message := fmt.Sprintf(
		"'%s' does not match the format. %s",
		name,
		messageDetails)

	return &UnprocessableError{
		&BaseError{
			Code:    "INCORRECT_GAME_NAME_FORMAT",
			Message: message,
		},
	}
}

func NewEmailUnprocessableError(email string, messageDetails string) error {
	message := fmt.Sprintf(
		"'%s' does not match the format. %s",
		email,
		messageDetails)

	return &UnprocessableError{
		&BaseError{
			Code:    "INCORRECT_EMAIL_FORMAT",
			Message: message,
		},
	}
}

func NewPasswordUnprocessableError(messageDetails string) error {
	message := fmt.Sprintf(
		"The password does not match the format. %s",
		messageDetails)

	return &UnprocessableError{
		&BaseError{
			Code:    "INCORRECT_PASSWORD_FORMAT",
			Message: message,
		},
	}
}

func NewChangeSourceUnprocessableError(possibleValues []string) error {
	message := fmt.Sprintf(
		"The change source doesn't match any of these: %s.",
		strings.Join(possibleValues, ", "))

	return &UnprocessableError{
		&BaseError{
			Code:    "INCORRECT_CHANGE_SOURCE_VALUE",
			Message: message,
		},
	}
}

func NewWheelEffectNameRequiredUnprocessableError(changeSource string) error {
	message := fmt.Sprintf(
		"A wheel effect name is required to do '%s' change source.",
		changeSource)

	return &UnprocessableError{
		&BaseError{
			Code:    "WHEEL_EFFECT_NAME_REQUIRED",
			Message: message,
		},
	}
}
