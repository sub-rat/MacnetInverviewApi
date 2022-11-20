package user

import "gorm.io/gorm"

type RepositoryInterface interface {
	Query(offset, limit int, query, fieldType string) ([]User, error)
	Get(id uint) (User, error)
	Create(req *User) error
	Update(id uint, update *User) error
	Delete(id uint) error
}

type repository struct {
	db gorm.DB
}

func NewRepository(db gorm.DB) RepositoryInterface {
	return &repository{db}
}

func (repository *repository) Query(offset, limit int, query, fieldType string) ([]User, error) {
	var dataList []User
	dbQuery := repository.db.Debug().Model(&User{})
	if fieldType == "email" {
		dbQuery.Where("email = ?", query)
	} else if fieldType == "username" {
		dbQuery.Where("username = ?", query)
	} else {
		dbQuery.Where("first_name like ? or user_name like ? or email like ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	err := dbQuery.Limit(limit).Offset(offset).
		Find(&dataList).
		Error
	return dataList, err
}

func (repository *repository) Get(id uint) (User, error) {
	user := User{}
	err := repository.db.Debug().
		Model(&User{}).
		First(&user, id).Error
	return user, err
}

func (repository *repository) Create(req *User) error {
	return repository.db.Debug().
		Model(&User{}).
		Create(&req).Error
}

func (repository *repository) Update(id uint, update *User) error {
	return repository.db.Debug().
		Model(&User{}).
		Where("id = ?", id).
		Updates(&update).Error
}

func (repository *repository) Delete(id uint) error {
	return repository.db.Debug().
		Delete(&User{}, id).Error
}
