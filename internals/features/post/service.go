package post

import (
	"errors"

	"github.com/sub-rat/social_network_api/internals/features/timeline"
	"github.com/sub-rat/social_network_api/internals/models"
)

type ServiceInterface interface {
	Query(offset, limit int, query string, fieldType string) ([]Post, error)
	Get(id uint) (Post, error)
	Create(req *Post) (Post, error)
	Update(id uint, update *Post) (Post, error)
	Delete(id uint) error
	PostLike(postLike *PostLike, post Post) error
	PostShare(postShare *PostShare, post Post) error
	PostShareDelete(postShare *PostShare) error
}

type service struct {
	repo         RepositoryInterface
	timelineRepo timeline.RepositoryInterface
}

type Post struct {
	models.Post
}

type PostLike struct {
	models.PostLike
}

type PostShare struct {
	models.PostShare
}

func NewService(repo RepositoryInterface, timelineRepo timeline.RepositoryInterface) ServiceInterface {
	return &service{repo, timelineRepo}
}

func (service *service) Query(offset, limit int, query string, fieldType string) ([]Post, error) {
	dataList, err := service.repo.Query(offset, limit, query, fieldType)
	if err != nil {
		return []Post{}, err
	}
	return dataList, nil
}

func (service *service) Get(id uint) (Post, error) {
	post, err := service.repo.Get(id)
	return post, err
}

func (service *service) Create(req *Post) (Post, error) {
	return *req, service.repo.Create(req)
}

func (service *service) Update(id uint, update *Post) (Post, error) {
	err := service.repo.Update(id, update)
	if err != nil {
		return Post{}, err
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

func (service *service) PostLike(postLike *PostLike, post Post) error {
	friend, _ := service.timelineRepo.FindFriend(postLike.UserId, post.UserId)
	if friend.RequestId != postLike.UserId {
		return errors.New("invalid post id")
	}
	pLike, _ := service.repo.GetPostShare(int(post.ID), postLike.UserId)
	if pLike.ID != 0 {
		return errors.New("post already liked")
	}
	return service.repo.PostLike(postLike)
}

func (service *service) PostShare(postShare *PostShare, post Post) error {
	friend, _ := service.timelineRepo.FindFriend(postShare.UserId, post.UserId)
	if friend.RequestId != postShare.UserId {
		return errors.New("invalid post id")
	}
	pShare, _ := service.repo.GetPostShare(int(post.ID), postShare.UserId)
	if pShare.ID != 0 {
		return errors.New("post already shared")
	}
	return service.repo.PostShare(postShare)
}

func (service *service) PostShareDelete(postShare *PostShare) error {
	return service.repo.PostShare(postShare)
}
