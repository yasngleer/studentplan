package store

import (
	"github.com/yasngleer/studentplan/models"
	"gorm.io/gorm"
)

type UserStore struct {
	DB *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	db.AutoMigrate(&models.User{})
	return &UserStore{DB: db}
}

func (s *UserStore) CreateUser(user *models.User) error {
	return s.DB.Create(user).Error
}

func (s *UserStore) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *UserStore) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	user.Username = username
	if err := s.DB.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) UpdateUser(user *models.User) error {
	return s.DB.Save(user).Error
}
