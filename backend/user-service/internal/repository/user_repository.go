package repository

import (
	"github.com/fastsprint/common/models"
	"github.com/fastsprint/common/utils"
	"github.com/fastsprint/user-service/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByUsernameOrEmail(usernameOrEmail string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	List(query model.UserQuery) ([]model.User, int64, error)
	GetByIDs(ids []uint) ([]model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsernameOrEmail(usernameOrEmail string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) List(query model.UserQuery) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := r.db.Model(&model.User{})

	if query.Keyword != "" {
		db = db.Where("username LIKE ? OR email LIKE ? OR display_name LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := utils.Paginate(query.Page, query.PageSize)
	err := db.Offset(offset).Limit(limit).Order("id desc").Find(&users).Error

	return users, total, err
}

func (r *userRepository) GetByIDs(ids []uint) ([]model.User, error) {
	var users []model.User
	if len(ids) == 0 {
		return users, nil
	}
	err := r.db.Where("id IN ?", ids).Find(&users).Error
	return users, err
}
