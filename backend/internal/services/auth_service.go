package services

import (
	"errors"

	"k-mate-backend/internal/dto"
	"k-mate-backend/internal/models"
	"k-mate-backend/internal/repository"
	"k-mate-backend/pkg/utils"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register creates a new user account
// Returns error if username already exists or password hashing fails
func (s *AuthService) Register(req *dto.RegisterRequest) (*models.User, error) {
	// 1. Check if username already exists
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, errors.New("gagal memeriksa username")
	}
	if exists {
		return nil, errors.New("username sudah digunakan")
	}

	// 2. Hash password using bcrypt
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}

	// 3. Create user model
	user := &models.User{
		FullName:     req.FullName,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Role:         models.UserRole(req.Role),
	}

	// 4. Save to database
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, errors.New("gagal menyimpan user ke database")
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
// Returns error if credentials are invalid
func (s *AuthService) Login(req *dto.LoginRequest) (string, *models.User, error) {
	// 1. Find user by username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return "", nil, errors.New("username atau password salah")
	}

	// 2. Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return "", nil, errors.New("username atau password salah")
	}

	// 3. Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, string(user.Role))
	if err != nil {
		return "", nil, errors.New("gagal membuat token autentikasi")
	}

	return token, user, nil
}
