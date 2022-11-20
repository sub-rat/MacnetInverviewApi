package timeline

import (
	"fmt"

	"gorm.io/gorm"
)

type RepositoryInterface interface {
	MyTimelines(offset, limit int, userId int) ([]Timeline, error)
	Dashboard(offset, limit int, userId int) ([]Timeline, error)
	Friends(offset, limit int, userId int, getRequested bool) ([]UserFriend, error)
	FindFriend(userId, friendId int) (UserFriend, error)
	AddFriend(userFriend *UserFriend) error
	UpdateFriend(iuserFriend *UserFriend) error
	RejectFriendRequest(userId, friendId int) error
	RemoveFriend(userId, friendId int) error
}

type repository struct {
	db gorm.DB
}

func NewRepository(db gorm.DB) RepositoryInterface {
	return &repository{db}
}

func (repository *repository) MyTimelines(offset, limit int, userId int) ([]Timeline, error) {
	var dataList []Timeline
	err := repository.db.Debug().Model(&Post{}).
		Select("posts.*,post_shares.id as shared_id").
		Joins("left join post_shares on posts.id = post_shares.post_id").
		Where("post_shares.user_id = ? or posts.user_id = ?", userId, userId).
		Limit(limit).Offset(offset).
		Find(&dataList).
		Error
	return dataList, err
}

func (repository *repository) Dashboard(offset, limit int, userId int) ([]Timeline, error) {
	var dataList []Timeline
	err := repository.db.Raw("? UNION ? limit ? offset ?",
		repository.db.Debug().Model(&Post{}).Select("posts.*,post_shares.id as shared_id").
			Joins("left join post_shares on posts.id = post_shares.post_id").
			Where("post_shares.user_id = ? or posts.user_id = ?", userId, userId),
		repository.db.Debug().Model(&Post{}).Select("posts.*,null as shared_id").
			Joins("left join user_friends on user_friends.accept_id = posts.user_id").
			Where("request_id = ? or posts.user_id = ? and is_accepted = true", userId, userId),
		limit, offset,
	).Scan(&dataList).Error
	return dataList, err
}

func (repository *repository) Friends(offset, limit int, userId int, getRequested bool) ([]UserFriend, error) {
	var dataList []UserFriend
	query := repository.db.Debug().Model(&UserFriend{})
	if getRequested {
		query.Preload("RequestUser").Where("accept_id = ?", userId)
	} else {
		query.Preload("AcceptUser").Where("request_id = ?", userId)
	}
	err := query.Limit(limit).Offset(offset).
		Find(&dataList).
		Error
	return dataList, err
}

func (repository *repository) FindFriend(userId, friendId int) (UserFriend, error) {
	friend := &UserFriend{}
	err := repository.db.Debug().
		Model(&UserFriend{}).
		Where("request_id = ? and accept_id = ?", userId, friendId).
		Find(friend).Limit(1).Error
	return *friend, err
}

func (repository *repository) AddFriend(userFriend *UserFriend) error {
	fmt.Println(userFriend)
	return repository.db.Debug().
		Model(&UserFriend{}).
		Create(&userFriend).Error
}

func (repository *repository) UpdateFriend(userFriend *UserFriend) error {
	fmt.Println(userFriend)
	return repository.db.Debug().
		Model(&UserFriend{}).
		Where("id = ?", userFriend.ID).
		UpdateColumn("is_accepted", userFriend.IsAccepted).Error
}

func (repository *repository) RejectFriendRequest(userId, friendId int) error {
	return repository.db.Debug().
		Where("request_id = ? and accept_id = ?", friendId, userId).
		Delete(&UserFriend{}).Error
}

func (repository *repository) RemoveFriend(userId, friendId int) error {
	return repository.db.Debug().
		Where("(request_id = ? and accept_id = ?) or (request_id = ? and accept_id = ?)", friendId, userId, userId, friendId).
		Delete(&UserFriend{}).Error
}
