package service

import (
	"crud/internal/model"
	"crud/internal/repository"
	"fmt"
)

const MaxUserLimit = 20

// IUserService is service layer for user so the handle layer can
// communicate with the datastore layer. User Service layer interface
// implementing business logic for user operation
type IUserService interface {
	Create(user *model.CreateUserRequest) (*model.UserResponse, error)
	GetById(id int) (*model.UserResponse, error)
	Update(user *model.UpdateUserRequest) (*model.UserResponse, error)
	Delete(id int) (*model.UserResponse, error)
	GetUsers(offset int, limit int) ([]*model.UserResponse, error)
}

// UserService is instance wrapper for IUserStore interface
type UserService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) IUserService {
	return &UserService{userRepository: userRepository}
}

func (service *UserService) Create(user *model.CreateUserRequest) (*model.UserResponse, error) {
	createUserModel := &model.UserModel{
		Name:  user.Name,
		Age:   user.Age,
		Email: user.Email,
	}
	newUserModel, err := service.userRepository.Create(createUserModel)
	if err != nil {
		return nil, err
	}
	return model.UserModelToUserResponse(newUserModel), nil
}

func (service *UserService) GetById(id int) (*model.UserResponse, error) {
	userModel, err := service.userRepository.GetById(id)
	if err != nil {
		return nil, err
	}
	return model.UserModelToUserResponse(userModel), nil
}

func (service *UserService) Update(user *model.UpdateUserRequest) (*model.UserResponse, error) {
	updateUserModel := &model.UserModel{
		ID:    user.Id,
		Name:  user.Name,
		Age:   user.Age,
		Email: user.Email,
	}
	userModel, err := service.userRepository.Update(updateUserModel)
	if err != nil {
		return nil, err
	}
	return model.UserModelToUserResponse(userModel), nil
}

func (service *UserService) Delete(id int) (*model.UserResponse, error) {
	userModel, err := service.userRepository.Delete(id)
	if err != nil {
		return nil, err
	}
	return model.UserModelToUserResponse(userModel), nil
}

func (service *UserService) GetUsers(offset int, limit int) ([]*model.UserResponse, error) {
	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be less than 0")
	}
	if limit > MaxUserLimit {
		return nil, fmt.Errorf("limit cannot be greater than %d", MaxUserLimit)
	} else if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than zero")
	}
	users, err := service.userRepository.GetAll(offset, limit)
	if err != nil {
		return nil, err
	}
	usersResponses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		usersResponses[i] = model.UserModelToUserResponse(user)
	}
	return usersResponses, nil
}
