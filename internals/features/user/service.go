package user

import (
	"errors"

	"github.com/sub-rat/social_network_api/internals/models"
	"github.com/sub-rat/social_network_api/pkg/utils"
)

type ServiceInterface interface {
	Query(offset, limit int, query string, fieldType string) ([]User, error)
	Get(id uint) (User, error)
	Create(req *User) (User, error)
	Update(id uint, update *User) (User, error)
	Delete(id uint) error
}

type service struct {
	repo RepositoryInterface
}

type User struct {
	models.User
}

func NewService(repo RepositoryInterface) ServiceInterface {
	return &service{repo}
}

func (service *service) Query(offset, limit int, query string, fieldType string) ([]User, error) {
	dataList, err := service.repo.Query(offset, limit, query, fieldType)
	if err != nil {
		return []User{}, err
	}
	return dataList, nil
}

func (service *service) Get(id uint) (User, error) {
	user, err := service.repo.Get(id)
	return user, err
}

func (service *service) Create(req *User) (User, error) {
	dataList, _ := service.repo.Query(0, 1, req.Email, "email")
	if len(dataList) > 0 {
		return User{}, errors.New("email already exists")
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return User{}, err
	}
	req.Password = hashedPassword
	return *req, service.repo.Create(req)
}

func (service *service) Update(id uint, update *User) (User, error) {
	err := service.repo.Update(id, update)
	if err != nil {
		return User{}, err
	}
	return *update, nil
}

func (service *service) Delete(id uint) error {
	err := service.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
