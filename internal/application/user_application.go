package application

import (
	"labelplus-next-web-be/internal/repository"
)

type UserApplication interface{}

type userApplication struct {
	userRepository repository.UserRepository
}

func NewUserApplication(
	userRepository repository.UserRepository,
) UserApplication {
	return &userApplication{
		userRepository: userRepository,
	}
}

// func (ua *userApplication) LoginUser(scope util.TraceScope) util.Result[value.LoginUserResponse] {
// }

// func (ua *userApplication) RegisterUser() util.Result[value.RegisterUserResponse] {
// }
