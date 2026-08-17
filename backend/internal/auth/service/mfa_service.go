package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"gorm.io/datatypes"

	"github.com/bees/hindu-ritual-platform/internal/auth/dto"
	"github.com/bees/hindu-ritual-platform/internal/auth/model"
	"github.com/bees/hindu-ritual-platform/internal/auth/repository"
	"github.com/bees/hindu-ritual-platform/pkg/configs"
	"github.com/bees/hindu-ritual-platform/pkg/mfa"
)

type MFAService struct {
	repo   *repository.AuthRepository
	config *configs.Config
}

func NewMFAService(repo *repository.AuthRepository, config *configs.Config) *MFAService {
	return &MFAService{
		repo:   repo,
		config: config,
	}
}

func (s *MFAService) SetupMFA(userID uuid.UUID) (*dto.MFASetupResponse, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	secret, err := mfa.GenerateTOTPSecret()
	if err != nil {
		return nil, errors.New("failed to generate MFA secret")
	}

	recoveryCodes := make([]string, 8)
	for i := range recoveryCodes {
		codes, _ := mfa.GenerateRecoveryCodes()
		recoveryCodes[i] = codes[0]
	}

	user.MFASecret = secret.Secret
	user.RecoveryCodes = datatypes.NewJSONSlice(recoveryCodes)

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, errors.New("failed to save MFA settings")
	}

	qrURL, err := mfa.GetTOTPQRCode(secret.Secret, user.Email)
	if err != nil {
		qrURL = ""
	}

	return &dto.MFASetupResponse{
		Secret:        secret.Secret,
		QRCodeURL:     qrURL,
		RecoveryCodes: recoveryCodes,
	}, nil
}

func (s *MFAService) VerifyMFA(userID uuid.UUID, code string) error {
	user, err := s.repo.FindUserByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	if !totp.Validate(code, user.MFASecret) {
		return errors.New("invalid MFA code")
	}

	user.MFAEnabled = true
	return s.repo.UpdateUser(user)
}

func (s *MFAService) ValidateMFALogin(user *model.User, code, recoveryCode string) error {
	if code != "" {
		if totp.Validate(code, user.MFASecret) {
			return nil
		}
	}

	if recoveryCode != "" {
		for i, rc := range user.RecoveryCodes {
			if rc == recoveryCode {
				user.RecoveryCodes = append(user.RecoveryCodes[:i], user.RecoveryCodes[i+1:]...)
				s.repo.UpdateUser(user)
				return nil
			}
		}
	}

	return errors.New("invalid MFA code or recovery code")
}

func (s *MFAService) DisableMFA(userID uuid.UUID) error {
	user, err := s.repo.FindUserByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	user.MFAEnabled = false
	user.MFASecret = ""
	user.RecoveryCodes = nil

	return s.repo.UpdateUser(user)
}
