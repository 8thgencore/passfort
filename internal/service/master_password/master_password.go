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
	user, err := svc.storage.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Error("failed to check master password", "error", err.Error())
		return false, domain.ErrInternal
	}

	return user.MasterPassword.Valid, nil
}

// SaveMasterPassword saves or updates the master password for the given user.
func (svc *MasterPasswordService) SaveMasterPassword(ctx context.Context, userID uuid.UUID, password string) error {
	// Retrieve the user based on the userID
	userDAO, err := svc.storage.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Error("failed get user", "error", err)
		return domain.ErrInternal
	}
	user := converter.ToUser(userDAO)

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		svc.log.Error("failed to hashing password", "error", err)
		return domain.ErrInternal
	}

	// Generation of a new salt
	salt, err := cipherkit.GenerateSalt()
	if err != nil {
		svc.log.Error("failed to generate salt", "error", err)
		return domain.ErrInternal
	}

	// Generating a new key
	newKey := cipherkit.DeriveKey(password, salt)

	// Store deriive key in cache with a TTL
	encryptionKey := util.GenerateCacheKey("encryption_key", userID.String())
	valueSerialized, err := util.Serialize(newKey)
	if err != nil {
		svc.log.Error("failed to serialize cache key", "error", err)
		return domain.ErrInternal
	}

	err = svc.cache.Set(ctx, encryptionKey, valueSerialized, svc.masterPasswordTTL)
	if err != nil {
		svc.log.Error("failed to store master password activation status", "error", err.Error())
		return domain.ErrInternal
	}

	user.MasterPassword = hashedPassword
	user.Salt = salt
	_, err = svc.storage.UpdateUser(ctx, converter.ToUserDAO(user))
	if err != nil {
		svc.log.Error("failed to update user", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// ChangeMasterPassword changes the master password for the given user.
func (svc *MasterPasswordService) ChangeMasterPassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	// Retrieve the user based on the userID
	userDAO, err := svc.storage.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	user := converter.ToUser(userDAO)

	// Verify the old master password
	err = util.CompareHash(oldPassword, user.MasterPassword)
	if err != nil {
		return domain.ErrInvalidMasterPassword
	}

	// Hash the new master password
	hashedNewPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return domain.ErrInternal
	}

	// Update the master password
	user.MasterPassword = hashedNewPassword
	_, err = svc.storage.UpdateUser(ctx, converter.ToUserDAO(user))
	if err != nil {
		svc.log.Error("failed to update user", "error", err.Error())
		return domain.ErrInternal
	}

	err = svc.secretSvc.ReencryptAllSecrets(ctx, userID, []byte(oldPassword), []byte(newPassword))
	if err != nil {
		svc.log.Error("failed to start reencrypt all secrets", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// ActivateMasterPassword activates the master password for the given user.
func (svc *MasterPasswordService) ActivateMasterPassword(ctx context.Context, userID uuid.UUID, password string) error {
	user, err := svc.storage.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Error("failed to activate master password", "error", err.Error())
		return domain.ErrInternal
	}

	if !user.MasterPassword.Valid || user.Salt == nil {
		return domain.ErrMasterPasswordNotSet
	}

	// Verify the master password
	err = util.CompareHash(password, user.MasterPassword.String)
	if err != nil {
		return domain.ErrInvalidMasterPassword
	}

	// Generating a new key
	newKey := cipherkit.DeriveKey(password, user.Salt)

	// Store deriive key in cache with a TTL
	encryptionKey := util.GenerateCacheKey("encryption_key", userID.String())
	valueSerialized, err := util.Serialize(newKey)
	if err != nil {
		svc.log.Error("failed to serialize cache key", "error", err)
		return domain.ErrInternal
	}

	err = svc.cache.Set(ctx, encryptionKey, valueSerialized, svc.masterPasswordTTL)
	if err != nil {
		svc.log.Error("failed to store master password activation status", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// GetEncryptionKey is required to encrypt or decrypt the password
func (svc *MasterPasswordService) GetEncryptionKey(ctx context.Context, userID uuid.UUID) ([]byte, error) {
	var encryptionKey []byte

	cacheKey := util.GenerateCacheKey("encryption_key", userID.String())
	valueSerialized, err := svc.cache.Get(ctx, cacheKey)
	if err != nil {
		svc.log.Error("failed to get encryption key", "error", err.Error())
		return nil, domain.ErrMasterPasswordActivationExpired
	}

	err = util.Deserialize(valueSerialized, &encryptionKey)
	if err != nil {
		svc.log.Error("failed to deserialize encryption key", "error", err.Error())
		return nil, domain.ErrInternal
	}

	return encryptionKey, nil
}
