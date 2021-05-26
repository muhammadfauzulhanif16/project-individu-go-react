package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"silih_a3/entities"
	"silih_a3/formatters"
	"silih_a3/inputs"
	"silih_a3/repositories"
)

type UserService interface {
	SignUpUser(signUpInput inputs.SignUpUserInput) (entities.User, error)
	SignInUser(signInInput inputs.SignInUserInput) (entities.User, error)
	IsEmailAvailable(emailInput inputs.CheckEmailInput) (bool, error)
	SaveAvatar(userId int, avatarPath string) (entities.User, error)
	GetAllUsers() ([]formatters.UserDataFormatter, error)
	GetUserById(userId int) (entities.User, error)
	GetAllUsersById(idInput inputs.UserIdInput) (entities.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *userService {
	return &userService{userRepository}
}

func (s *userService) SignUpUser(signUpInput inputs.SignUpUserInput) (entities.User, error) {
	userData := entities.User{}
	userData.FirstName = signUpInput.FirstName
	userData.LastName = signUpInput.LastName
	userData.Email = signUpInput.Email
	passwordHashData, err := bcrypt.GenerateFromPassword([]byte(signUpInput.Password), bcrypt.MinCost)
	if err != nil {
		return userData, err
	}
	userData.Password = string(passwordHashData)
	userData.Role = "user"

	newUser, err := s.userRepository.InsertUser(userData)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *userService) SignInUser(signInInput inputs.SignInUserInput) (entities.User, error) {
	userEmail := signInInput.Email
	userPassword := signInInput.Password
	userData, err := s.userRepository.FindUserByEmail(userEmail)
	if err != nil {
		return userData, err
	}

	if userData.Id == 0 {
		return userData, errors.New("no user found with this email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(userPassword))
	if err != nil {
		return userData, err
	}

	return userData, nil
}

func (s *userService) IsEmailAvailable(emailInput inputs.CheckEmailInput) (bool, error) {
	emailData := emailInput.Email
	userData, err := s.userRepository.FindUserByEmail(emailData)
	if err != nil {
		return false, err
	}

	if userData.Id == 0 {
		return true, nil
	}

	return false, nil
}

func (s *userService) SaveAvatar(userId int, avatarPath string) (entities.User, error) {
	userData, err := s.userRepository.FindUserById(userId)
	if err != nil {
		return userData, err
	}

	userData.Avatar = avatarPath
	updateUser, err := s.userRepository.UpdateUser(userData)
	if err != nil {
		return updateUser, err
	}

	return updateUser, nil
}

func (s *userService) GetAllUsers() ([]formatters.UserDataFormatter, error) {
	var userDataFormat []formatters.UserDataFormatter
	users, err := s.userRepository.FindAllUsers()
	for _, user := range users {
		userDataFormat = append(userDataFormat, formatters.UserFormat(user))
	}

	if err != nil {
		return userDataFormat, err
	}

	return userDataFormat, nil
}

func (s *userService) GetUserById(userId int) (entities.User, error) {
	userData, err := s.userRepository.FindUserById(userId)
	if err != nil {
		return userData, err
	}

	if userData.Id == 0 {
		return userData, errors.New("there aren't users with this id")
	}

	return userData, nil
}

func (s *userService) GetAllUsersById(idInput inputs.UserIdInput) (entities.User, error) {
	userData, err := s.userRepository.FindUserById(idInput.Id)
	if err != nil {
		return userData, err
	}

	return userData, nil
}
