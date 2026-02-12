package validator

import (
	"FGG-Service/src/common"
	"regexp"
)

var userNameRegex = regexp.MustCompile(`^\w+$`)

func ValidateUserName(name string) error {
	if len(name) < 3 {
		return common.NewUserNameBadRequestError(
			name,
			"The name is too short, it should be at least 3 characters long.")
	}

	if len(name) > 35 {
		return common.NewUserNameBadRequestError(
			name,
			"The name is too long, it should be less than 35 characters long.")
	}

	if !userNameRegex.MatchString(name) {
		return common.NewUserNameBadRequestError(
			name,
			"The name must contain only Latin letters, numbers, and underscores")
	}

	return nil
}

var gameNameRegex = regexp.MustCompile(`^.+$`)

func ValidateGameName(name string) error {
	if len(name) < 1 {
		return common.NewGameNameBadRequestError(
			name,
			"The name is too short, it should be at least 1 characters long.")
	}

	if len(name) > 70 {
		return common.NewGameNameBadRequestError(
			name,
			"The name is too long, it should be less than 70 characters long.")
	}

	if !gameNameRegex.MatchString(name) {
		return common.NewGameNameBadRequestError(
			name,
			"The name must contain only UTF-8 characters")
	}

	return nil
}

var emailRegex = regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)

func ValidateEmail(email string) error {
	if len(email) < 6 {
		return common.NewEmailBadRequestError(
			email,
			"The email is too short, it should be at least 6 characters long.")
	}

	if len(email) > 100 {
		return common.NewEmailBadRequestError(
			email,
			"The email is too long, it should be less than 100 characters long.")
	}

	if !emailRegex.MatchString(email) {
		return common.NewEmailBadRequestError(
			email,
			"The email must follow the pattern (e.g. email@example.com)")
	}

	return nil
}

var passwordUppercaseLetterRegex = regexp.MustCompile(`[A-Z]`)
var passwordLowercaseLetterRegex = regexp.MustCompile(`[a-z]`)
var passwordDigitRegex = regexp.MustCompile(`\d`)
var passwordSpecialSymbolRegex = regexp.MustCompile(`[ !"#$%&'()*+,-./:;<=>?@\[\\\]^_{|}~]`)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return common.NewPasswordBadRequestError(
			"The password is too short, it should be at least 8 characters long.")
	}

	if len(password) > 35 {
		return common.NewPasswordBadRequestError(
			"The password is too long, it should be less than 35 characters long.")
	}

	if !passwordUppercaseLetterRegex.MatchString(password) {
		return common.NewPasswordBadRequestError(
			"The password doesn't contain uppercase letters.")
	}

	if !passwordLowercaseLetterRegex.MatchString(password) {
		return common.NewPasswordBadRequestError(
			"The password doesn't contain lowercase letters.")
	}

	if !passwordDigitRegex.MatchString(password) {
		return common.NewPasswordBadRequestError(
			"The password doesn't contain digits.")
	}

	if !passwordSpecialSymbolRegex.MatchString(password) {
		return common.NewPasswordBadRequestError(
			"The password doesn't contain any of special symbols ( !\"#$%&'()*+,-./:;<=>?@[]\\^_{|}~).")
	}

	return nil
}
