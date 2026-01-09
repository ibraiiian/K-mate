package repository

import (
	"k-mate-backend/config"
	"k-mate-backend/internal/models"
)

// UserRepository handles database operations for User model
type UserRepository struct{}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

// FindByUsername retrieves a user by their username
// Returns nil if user not found
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	result := config.DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// FindByID retrieves a user by their ID
func (r *UserRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	result := config.DB.First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// ExistsByUsername checks if a username already exists
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := config.DB.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
