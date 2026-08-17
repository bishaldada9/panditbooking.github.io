package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/bees/hindu-ritual-platform/internal/auth/model"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *AuthRepository) FindUserByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *AuthRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *AuthRepository) IncrementFailedAttempts(userID uuid.UUID) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("failed_login_attempts", gorm.Expr("failed_login_attempts + 1")).Error
}

func (r *AuthRepository) ResetFailedAttempts(userID uuid.UUID) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"failed_login_attempts": 0,
			"locked_until":          nil,
		}).Error
}

func (r *AuthRepository) LockAccount(userID uuid.UUID) error {
	lockDuration := 15 * time.Minute
	lockTime := time.Now().Add(lockDuration)
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Update("locked_until", lockTime).Error
}

func (r *AuthRepository) UpdateLastLogin(userID uuid.UUID, ip, userAgent string) error {
	now := time.Now()
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at":   &now,
			"last_login_ip":   ip,
			"last_user_agent": userAgent,
		}).Error
}

func (r *AuthRepository) CreateRefreshToken(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) FindRefreshToken(token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.Where("token = ? AND is_revoked = ? AND expires_at > ?", token, false, time.Now()).First(&rt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rt, err
}

func (r *AuthRepository) RevokeRefreshToken(id uuid.UUID) error {
	return r.db.Model(&model.RefreshToken{}).Where("id = ?", id).
		Update("is_revoked", true).Error
}

func (r *AuthRepository) RevokeAllUserTokens(userID uuid.UUID) error {
	return r.db.Model(&model.RefreshToken{}).Where("user_id = ? AND is_revoked = ?", userID, false).
		Update("is_revoked", true).Error
}

func (r *AuthRepository) CreateOTP(otp *model.OTP) error {
	return r.db.Create(otp).Error
}

func (r *AuthRepository) FindValidOTP(userID uuid.UUID, code, purpose string) (*model.OTP, error) {
	var otp model.OTP
	query := r.db.Where("code = ? AND purpose = ? AND is_used = ? AND expires_at > ?",
		code, purpose, false, time.Now())
	if userID != uuid.Nil {
		query = query.Where("user_id = ?", userID)
	}
	err := query.First(&otp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &otp, err
}

func (r *AuthRepository) MarkOTPUsed(id uuid.UUID) error {
	return r.db.Model(&model.OTP{}).Where("id = ?", id).
		Update("is_used", true).Error
}

func (r *AuthRepository) UpdatePassword(userID uuid.UUID, passwordHash string) error {
	now := time.Now()
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password_hash":       passwordHash,
			"password_changed_at": &now,
		}).Error
}

func (r *AuthRepository) VerifyEmail(userID uuid.UUID) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Update("is_email_verified", true).Error
}

func (r *AuthRepository) GetUserSessions(userID uuid.UUID) ([]model.RefreshToken, error) {
	var tokens []model.RefreshToken
	err := r.db.Where("user_id = ? AND is_revoked = ? AND expires_at > ?",
		userID, false, time.Now()).Find(&tokens).Error
	return tokens, err
}
