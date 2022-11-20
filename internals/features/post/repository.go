package post

import "gorm.io/gorm"

type RepositoryInterface interface {
	Query(offset, limit int, query, fieldType string) ([]Post, error)
	Get(id uint) (Post, error)
	Create(req *Post) error
	Update(id uint, update *Post) error
	Delete(id uint) error
	PostLike(postLike *PostLike) error
	PostShare(postShare *PostShare) error
	PostShareDelete(postShare *PostShare) error
	GetPostShare(postId int, userId int) (PostShare, error)
	GetPostLike(postId int, userId int) (PostLike, error)
}

type repository struct {
	db gorm.DB
}

func NewRepository(db gorm.DB) RepositoryInterface {
	return &repository{db}
}

func (repository *repository) Query(offset, limit int, query, fieldType string) ([]Post, error) {
	var dataList []Post
	dbQuery := repository.db.Debug().Model(&Post{}).Preload("PostLike").Preload("PostShare")
	if fieldType == "message" {
		dbQuery.Where("message = ?", query)
	} else if fieldType == "user_id" {
		dbQuery.Where("user_id = ?", query)
	} else {
		dbQuery.Where("message like ? ", "%"+query+"%")
	}

	err := dbQuery.Limit(limit).Offset(offset).
		Find(&dataList).
		Error
	return dataList, err
}

func (repository *repository) Get(id uint) (Post, error) {
	post := Post{}
	err := repository.db.Debug().
		Model(&Post{}).
		Preload("PostLike").
		Preload("PostShare").
		Preload("User").
		First(&post, id).Error
	return post, err
}

func (repository *repository) Create(req *Post) error {
	return repository.db.Debug().
		Model(&Post{}).
		Create(&req).Error
}

func (repository *repository) Update(id uint, update *Post) error {
	return repository.db.Debug().
		Model(&Post{}).
		Where("id = ?", id).
		Updates(&update).Error
}

func (repository *repository) Delete(id uint) error {
	return repository.db.Debug().
		Delete(&Post{}, id).Error
}

func (repository *repository) PostLike(postLike *PostLike) error {
	plike := &PostLike{}
	err := repository.db.Debug().
		Model(&PostLike{}).
		Find(plike).Error
	if err == nil && plike.ID == 0 {
		return repository.db.Debug().
			Model(&PostLike{}).
			Create(&postLike).Error
	}
	return repository.db.Debug().
		Model(&PostLike{}).
		Delete(plike).Error
}

func (repository *repository) PostShare(postShare *PostShare) error {
	return repository.db.Debug().
		Model(&PostShare{}).
		Create(&postShare).Error

}

func (repository *repository) PostShareDelete(postShare *PostShare) error {
	return repository.db.Debug().
		Model(&PostShare{}).
		Delete(postShare).Error
}

func (repository *repository) GetPostShare(postId int, userId int) (PostShare, error) {
	postShare := PostShare{}
	err := repository.db.Debug().
		Model(&PostShare{}).
		Where("user_id = ? and post_id = ?", userId, postId).
		First(&postShare).Error
	return postShare, err
}

func (repository *repository) GetPostLike(postId int, userId int) (PostLike, error) {
	postLike := PostLike{}
	err := repository.db.Debug().
		Model(&PostLike{}).
		Where("user_id = ? and post_id = ?", userId, postId).
		First(&postLike).Error
	return postLike, err
}
