package user

import (
	"context"
	"errors"
	"funding-app/app/helper"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	LoginUser(input LoginUserInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	UploadAvatar(userID string, fileLocation string) (User, error)
	GetUserByID(userID string) (User, error)
}

type service struct {
	userRepository Repository
}

func NewService(userRepository Repository) Service {
	return &service{userRepository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	var user User
	userID := helper.GenerateID()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user.ID = userID
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	password := string(passwordHash)
	user.PasswordHash = password
	user.Role = "user"

	newUser, err := s.userRepository.Save(ctx, user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *service) LoginUser(input LoginUserInput) (User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user, err := s.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return user, err
	}

	if user.ID == "" {
		return user, errors.New("no user found on that email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return user, err
	}

	return user, err
}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user, err := s.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return false, err
	}

	if user.ID == "" {
		return true, nil
	}

	return false, nil
}

func (s *service) UploadAvatar(userID string, fileLocation string) (User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	updatedUser, err := s.userRepository.Update(ctx, user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *service) GetUserByID(userID string) (User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return user, err
	}

	if user.ID == "" {
		return user, errors.New("no user found")
	}

	return user, nil
}
