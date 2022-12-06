package timeline

import (
	"errors"

	"github.com/sub-rat/social_network_api/internals/features/user"
	"github.com/sub-rat/social_network_api/internals/models"
)

type ServiceInterface interface {
	MyTimeline(offset, limit int, userId int) ([]Timeline, error)
	Friends(offset, limit int, userId int, getRequested bool) ([]UserFriend, error)
	Dashboard(offset, limit int, userId int) ([]Timeline, error)
	AddFriend(userFriend *UserFriend) error
	AcceptFriendRequest(userId, friendId int) error
	RejectFriendRequest(userId, friendId int) error
	RemoveFriend(userId, friendId int) error
}

type service struct {
	repo     RepositoryInterface
	userRepo user.RepositoryInterface
}

type User struct {
	models.User
}

type UserFriend struct {
	models.UserFriend
}

type Post struct {
	models.Post
}

type PostShare struct {
	models.PostShare
}

type PostLike struct {
	models.PostLike
}

type Timeline struct {
	models.Post
	SharedId int
}

func NewService(repo RepositoryInterface, userRepo user.RepositoryInterface) ServiceInterface {
	return &service{repo, userRepo}
}

func (service *service) MyTimeline(offset, limit int, userId int) ([]Timeline, error) {
	dataList, err := service.repo.MyTimelines(offset, limit, userId)
	if err != nil {
		return []Timeline{}, err
	}
	return dataList, nil
}

func (service *service) Dashboard(offset, limit int, userId int) ([]Timeline, error) {
	dataList, err := service.repo.Dashboard(offset, limit, userId)
	if err != nil {
		return []Timeline{}, err
	}
	return dataList, nil
}

func (service *service) Friends(offset, limit int, userId int, getRequested bool) ([]UserFriend, error) {
	dataList, err := service.repo.Friends(offset, limit, userId, getRequested)
	if err != nil {
		return []UserFriend{}, err
	}
	return dataList, nil
}

func (service *service) AddFriend(userFriend *UserFriend) error {
	friend, err := service.repo.FindFriend(userFriend.RequestId, userFriend.AcceptId)
	if err == nil && friend.ID > 0 {
		return errors.New("friend request already send")
	}
	user, err := service.userRepo.Get(uint(userFriend.AcceptId))
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errors.New("friend doesn't exists")
	}
	return service.repo.AddFriend(userFriend)
}

// 1//3
func (service *service) AcceptFriendRequest(userId, friendId int) error {
	d1, _ := service.repo.FindFriend(userId, friendId)

	d2, _ := service.repo.FindFriend(friendId, userId)
	if d1.AcceptId == friendId && d2.AcceptId == userId {
		return errors.New(" request has been already accepted")
	} else if d2.AcceptId != userId {
		return errors.New("no request found")
	}

	d2.IsAccepted = true
	service.repo.UpdateFriend(&d2)
	tmp := d2.AcceptId
	d2.ID = 0
	d2.AcceptId = d2.RequestId
	d2.RequestId = tmp
	return service.repo.AddFriend(&d2)
}

func (service *service) RejectFriendRequest(userId, friendId int) error {
	_, err := service.repo.FindFriend(friendId, userId)
	if err != nil {
		return errors.New("no request found")
	}
	return service.repo.RejectFriendRequest(userId, friendId)
}

func (service *service) RemoveFriend(userId, friendId int) error {
	_, err := service.repo.FindFriend(userId, friendId)
	if err != nil {
		return errors.New("no friend found")
	}
	return service.repo.RemoveFriend(userId, friendId)
}
