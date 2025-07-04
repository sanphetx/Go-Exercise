package repository

import (
	"Go-Exercise/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
	FindUserByID(id uuid.UUID) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uuid.UUID) error

	SaveRefreshToken(token *model.RefreshToken) error
	GetRefreshToken(token string) (*model.RefreshToken, error)
	DeleteRefreshToken(token string) error
	UpdateRefreshToken(token *model.RefreshToken) error
	RevokeAllUserTokens(userID uuid.UUID) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindUserByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) SaveRefreshToken(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *authRepository) GetRefreshToken(token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := r.db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *authRepository) UpdateRefreshToken(token *model.RefreshToken) error {
	return r.db.Save(token).Error
}

func (r *authRepository) RevokeAllUserTokens(userID uuid.UUID) error {
	return r.db.Model(&model.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error
}

func (r *authRepository) DeleteRefreshToken(token string) error {
	return r.db.Delete(&model.RefreshToken{}, "token = ?", token).Error
}

func (r *authRepository) GetAllUsers() ([]*model.User, error) {
	var users []*model.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *authRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *authRepository) DeleteUser(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.RefreshToken{}, "user_id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Delete(&model.User{}, "id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}
