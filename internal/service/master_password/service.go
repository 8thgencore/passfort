package masterpassword

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/8thgencore/passfort/pkg/util"
	"github.com/google/uuid"
)

// MasterPasswordExists checks if a master password already exists for the given user.
func (s *MasterPasswordService) MasterPasswordExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		s.log.Error("failed to check master password", "error", err.Error())
		return false, domain.ErrInternal
	}

	return user.MasterPassword.Valid, nil
}

// SaveMasterPassword saves or updates the master password for the given user.
func (s *MasterPasswordService) SaveMasterPassword(ctx context.Context, userID uuid.UUID, password string) error {
	// Retrieve the user based on the userID
	userDAO, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	user := converter.ToUser(userDAO)

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return domain.ErrInternal
	}

	user.MasterPassword = hashedPassword
	_, err = s.storage.UpdateUser(ctx, converter.ToUserDAO(user))
	if err != nil {
		s.log.Error("failed to update user", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// ActivateMasterPassword activates the master password for the given user.
func (s *MasterPasswordService) ActivateMasterPassword(ctx context.Context, userID uuid.UUID, password string) error {
	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		s.log.Error("failed to activate master password", "error", err.Error())
		return domain.ErrInternal
	}

	if !user.MasterPassword.Valid {
		return domain.ErrMasterPasswordNotSet
	}

	err = util.ComparePassword(password, user.MasterPassword.String)
	if err != nil {
		return domain.ErrInvalidMasterPassword

	}

	// Store activation status in cache with a TTL
	cacheKey := util.GenerateCacheKey("master_password_activated", userID.String())
	valueSerialized, err := util.Serialize(user)
	if err != nil {
		return domain.ErrInternal
	}

	err = s.cache.Set(ctx, cacheKey, valueSerialized, s.masterPasswordTTL)
	if err != nil {
		s.log.Error("failed to store master password activation status", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// IsMasterPasswordActivated checks if the master password has been activated recently.
func (s *MasterPasswordService) IsMasterPasswordActivated(ctx context.Context, userID uuid.UUID) (bool, error) {
	var activated bool

	cacheKey := util.GenerateCacheKey("master_password_activated", userID.String())
	valueSerialized, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		s.log.Error("failed to get master password activation status", "error", err.Error())
		return false, domain.ErrMasterPasswordActivationExpired
	}

	err = util.Deserialize(valueSerialized, &activated)
	if err != nil {
		return false, domain.ErrInternal
	}

	return activated, nil
}
