package masterpassword

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/8thgencore/passfort/pkg/cipherkit"
	"github.com/8thgencore/passfort/pkg/util"
	"github.com/google/uuid"
)

// MasterPasswordExists checks if a master password already exists for the given user.
func (svc *MasterPasswordService) MasterPasswordExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	user, err := svc.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Error("Failed to check master password", "error", err.Error())
		return false, domain.ErrInternal
	}

	return user.MasterPassword.Valid, nil
}

// SaveMasterPassword saves or updates the master password for the given user.
func (svc *MasterPasswordService) SaveMasterPassword(ctx context.Context, userID uuid.UUID, password string) error {
	userDAO, err := svc.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Error("Failed to get user", "error", err.Error())
		return domain.ErrInternal
	}

	user := converter.ToUser(userDAO)

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		svc.log.Error("Failed to hash password", "error", err.Error())
		return domain.ErrInternal
	}

	salt, err := cipherkit.GenerateSalt()
	if err != nil {
		svc.log.Error("Failed to generate salt", "error", err.Error())
		return domain.ErrInternal
	}

	newKey := cipherkit.DeriveKey(password, salt)

	encryptionKey := util.GenerateCacheKey("encryption_key", userID.String())
	valueSerialized, err := util.Serialize(newKey)
	if err != nil {
		svc.log.Error("Failed to serialize cache key", "error", err.Error())
		return domain.ErrInternal
	}

	if err = svc.cache.Set(ctx, encryptionKey, valueSerialized, svc.masterPasswordTTL); err != nil {
		svc.log.Error("Failed to store master password activation status", "error", err.Error())
		return domain.ErrInternal
	}

	user.MasterPassword = hashedPassword
	user.Salt = salt

	updatedUser, err := svc.userStorage.UpdateUser(ctx, converter.ToUserDAO(user))
	if err != nil {
		svc.log.Error("Failed to update user", "error", err.Error())
		return domain.ErrInternal
	}

	cacheKey := util.GenerateCacheKey("user", updatedUser.ID)
	if err = svc.cache.Delete(ctx, cacheKey); err != nil {
		return domain.ErrInternal
	}

	serializedUser, err := util.Serialize(converter.ToUser(updatedUser))
	if err != nil {
		return domain.ErrInternal
	}

	if err = svc.cache.Set(ctx, cacheKey, serializedUser, 0); err != nil {
		return domain.ErrInternal
	}

	return nil
}

// ChangeMasterPassword changes the master password for the given user.
func (svc *MasterPasswordService) ChangeMasterPassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	userDAO, err := svc.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		return domain.ErrDataNotFound
	}

	user := converter.ToUser(userDAO)

	if err = util.CompareHash(oldPassword, user.MasterPassword); err != nil {
		return domain.ErrInvalidMasterPassword
	}

	oldEncryptionKey := cipherkit.DeriveKey(oldPassword, user.Salt)
	newEncryptionKey := cipherkit.DeriveKey(newPassword, user.Salt)

	hashedNewPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return domain.ErrInternal
	}

	user.MasterPassword = hashedNewPassword

	_, err = svc.userStorage.UpdateUser(ctx, converter.ToUserDAO(user))
	if err != nil {
		svc.log.Error("Failed to update user", "error", err.Error())
		return domain.ErrInternal
	}

	if err = svc.secretSvc.ReencryptAllSecrets(ctx, userID, oldEncryptionKey, newEncryptionKey); err != nil {
		svc.log.Error("Failed to start re-encrypt all secrets", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// ActivateMasterPassword activates the master password for the given user.
func (svc *MasterPasswordService) ActivateMasterPassword(ctx context.Context, userID uuid.UUID, password string) error {
	userDAO, err := svc.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Error("Failed to activate master password", "error", err.Error())
		return domain.ErrInternal
	}

	if !userDAO.MasterPassword.Valid || userDAO.Salt == nil {
		return domain.ErrMasterPasswordNotSet
	}

	if err = util.CompareHash(password, userDAO.MasterPassword.String); err != nil {
		return domain.ErrInvalidMasterPassword
	}

	newKey := cipherkit.DeriveKey(password, userDAO.Salt)

	encryptionKey := util.GenerateCacheKey("encryption_key", userID.String())
	valueSerialized, err := util.Serialize(newKey)
	if err != nil {
		svc.log.Error("Failed to serialize cache key", "error", err.Error())
		return domain.ErrInternal
	}

	if err = svc.cache.Set(ctx, encryptionKey, valueSerialized, svc.masterPasswordTTL); err != nil {
		svc.log.Error("Failed to store master password activation status", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// GetEncryptionKey retrieves the encryption key from the cache.
func (svc *MasterPasswordService) GetEncryptionKey(ctx context.Context, userID uuid.UUID) ([]byte, error) {
	cacheKey := util.GenerateCacheKey("encryption_key", userID.String())
	valueSerialized, err := svc.cache.Get(ctx, cacheKey)
	if err != nil {
		svc.log.Error("Failed to get encryption key", "error", err.Error())
		return nil, domain.ErrMasterPasswordActivationExpired
	}

	var encryptionKey []byte
	if err = util.Deserialize(valueSerialized, &encryptionKey); err != nil {
		svc.log.Error("Failed to deserialize encryption key", "error", err.Error())
		return nil, domain.ErrInternal
	}

	return encryptionKey, nil
}
