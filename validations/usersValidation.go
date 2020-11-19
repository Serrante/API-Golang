package validations

import (
	"api/models"
	"errors"
)

var (
	ErrEmptyFields  = errors.New("Um ou mais campos não foram preenchidos")
	ErrInvalidEmail = errors.New("E-mail inválido")
)

func ValidateNewUser(user models.User) (models.User, error) {
	if IsEmpty(user.Nickname) || IsEmpty(user.Email) || IsEmpty(user.Password) {
		return models.User{}, ErrEmptyFields
	}

	if !IsMail(user.Email) {
		return models.User{}, ErrInvalidEmail
	}

	return user, nil
}
