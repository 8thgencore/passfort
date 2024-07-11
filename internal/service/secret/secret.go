package secret

import (
	"context"
	"time"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/dao"
	"github.com/8thgencore/passfort/pkg/cipherkit"
	"github.com/8thgencore/passfort/pkg/logger/sl"
	"github.com/google/uuid"
)

// CreateSecret creates a new secret
func (svc *SecretService) CreateSecret(ctx context.Context, userID uuid.UUID, secret *domain.Secret, encryptionKey []byte) (*domain.Secret, error) {
	if !svc.isUserPartOfCollection(ctx, userID, secret.CollectionID) {
		return nil, domain.ErrUnauthorized
	}

	secret.CreatedBy = userID
	secret.UpdatedBy = userID
	secretDAO := converter.ToSecretDAO(secret)

	var err error
	switch secret.SecretType {
	case domain.PasswordSecretType:
		err = svc.createPasswordSecret(ctx, secret, encryptionKey, secretDAO)
	case domain.TextSecretType:
		err = svc.createTextSecret(ctx, secret, encryptionKey, secretDAO)
	default:
		return nil, domain.ErrInvalidSecretType
	}

	if err != nil {
		return nil, err
	}

	createdSecretDAO, err := svc.secretStorage.CreateSecret(ctx, secret.CollectionID, secretDAO)
	if err != nil {
		svc.log.Error("Error creating secret:", sl.Err(err))
		return nil, domain.ErrInternal
	}

	createdSecretDAO.TextSecret = secretDAO.TextSecret
	createdSecretDAO.PasswordSecret = secretDAO.PasswordSecret

	return converter.ToSecret(createdSecretDAO), nil
}

func (svc *SecretService) createPasswordSecret(ctx context.Context, secret *domain.Secret, encryptionKey []byte, secretDAO *dao.SecretDAO) error {
	encryptedPassword, err := cipherkit.Encrypt([]byte(secret.PasswordSecret.Password), encryptionKey)
	if err != nil {
		svc.log.Error("Error encrypting password secret:", sl.Err(err))
		return domain.ErrInternal
	}

	passwordSecretDAO := converter.ToPasswordSecretDAO(secret.PasswordSecret)
	passwordSecretDAO.Password = encryptedPassword

	newSecret, err := svc.secretStorage.CreatePasswordSecret(ctx, passwordSecretDAO)
	if err != nil {
		svc.log.Error("Error creating password secret:", sl.Err(err))
		return domain.ErrInternal
	}

	secretDAO.LinkedSecretId = newSecret.ID
	return nil
}

func (svc *SecretService) createTextSecret(ctx context.Context, secret *domain.Secret, encryptionKey []byte, secretDAO *dao.SecretDAO) error {
	encryptedText, err := cipherkit.Encrypt([]byte(secret.TextSecret.Text), encryptionKey)
	if err != nil {
		svc.log.Error("Error encrypting text secret:", sl.Err(err))
		return domain.ErrInternal
	}

	textSecretDAO := converter.ToTextSecretDAO(secret.TextSecret)
	textSecretDAO.Text = encryptedText

	newSecret, err := svc.secretStorage.CreateTextSecret(ctx, textSecretDAO)
	if err != nil {
		svc.log.Error("Error creating text secret:", sl.Err(err))
		return domain.ErrInternal
	}
	secretDAO.LinkedSecretId = newSecret.ID
	return nil
}

// ListSecretsByCollectionID lists secrets for a specific collection ID
func (svc *SecretService) ListSecretsByCollectionID(ctx context.Context, userID, collectionID uuid.UUID, skip, limit uint64) ([]domain.Secret, error) {
	if !svc.isUserPartOfCollection(ctx, userID, collectionID) {
		return nil, domain.ErrUnauthorized
	}

	secretsDAO, err := svc.secretStorage.ListSecretsByCollectionID(ctx, collectionID, skip, limit)
	if err != nil {
		svc.log.Error("Error listing secrets for collection:", "collectionID", collectionID, sl.Err(err))
		return nil, domain.ErrDataNotFound
	}

	secrets := make([]domain.Secret, 0, len(secretsDAO))
	for _, secretDAO := range secretsDAO {
		secrets = append(secrets, *converter.ToSecret(&secretDAO))
	}

	return secrets, nil
}

// GetSecret gets a secret by ID
func (svc *SecretService) GetSecret(ctx context.Context, userID, collectionID, secretID uuid.UUID, encryptionKey []byte) (*domain.Secret, error) {
	if !svc.isUserPartOfCollection(ctx, userID, collectionID) {
		return nil, domain.ErrUnauthorized
	}

	secretDAO, err := svc.secretStorage.GetSecretByID(ctx, secretID)
	if err != nil {
		svc.log.Error("Error getting secret by ID", "secretID", secretID, "error", err)
		return nil, domain.ErrDataNotFound
	}

	switch secretDAO.SecretType {
	case dao.PasswordSecretType:
		passwordSecretDAO, err := svc.getAndDecryptPasswordSecret(ctx, secretDAO.LinkedSecretId, encryptionKey)
		if err != nil {
			return nil, err
		}
		secretDAO.PasswordSecret = *passwordSecretDAO
	case dao.TextSecretType:
		textSecretDAO, err := svc.getAndDecryptTextSecret(ctx, secretDAO.LinkedSecretId, encryptionKey)
		if err != nil {
			return nil, err
		}
		secretDAO.TextSecret = *textSecretDAO
	default:
		svc.log.Error("Invalid secret type", "secretType", secretDAO.SecretType, "secretID", secretID)
		return nil, domain.ErrInvalidSecretType
	}

	return converter.ToSecret(secretDAO), nil
}

func (svc *SecretService) getAndDecryptPasswordSecret(ctx context.Context, linkedSecretID uuid.UUID, encryptionKey []byte) (*dao.PasswordSecretDAO, error) {
	passwordSecretDAO, err := svc.secretStorage.GetPasswordSecretByID(ctx, linkedSecretID)
	if err != nil {
		svc.log.Error("Error getting password secret by ID", "linkedSecretID", linkedSecretID, "error", err)
		return nil, domain.ErrDataNotFound
	}

	decryptedPassword, err := cipherkit.Decrypt(passwordSecretDAO.Password, encryptionKey)
	if err != nil {
		svc.log.Error("Error decrypting password secret", "linkedSecretID", linkedSecretID, "error", err)
		return nil, domain.ErrInternal
	}

	passwordSecretDAO.Password = decryptedPassword
	return passwordSecretDAO, nil
}

func (svc *SecretService) getAndDecryptTextSecret(ctx context.Context, linkedSecretID uuid.UUID, encryptionKey []byte) (*dao.TextSecretDAO, error) {
	textSecretDAO, err := svc.secretStorage.GetTextSecretByID(ctx, linkedSecretID)
	if err != nil {
		svc.log.Error("Error getting text secret by ID", "linkedSecretID", linkedSecretID, "error", err)
		return nil, domain.ErrDataNotFound
	}

	decryptedText, err := cipherkit.Decrypt(textSecretDAO.Text, encryptionKey)
	if err != nil {
		svc.log.Error("Error decrypting text secret", "linkedSecretID", linkedSecretID, "error", err)
		return nil, domain.ErrInternal
	}

	textSecretDAO.Text = decryptedText
	return textSecretDAO, nil
}

// UpdateSecret updates a secret
func (svc *SecretService) UpdateSecret(ctx context.Context, userID, collectionID uuid.UUID, secret *domain.Secret, encryptionKey []byte) (*domain.Secret, error) {
	if !svc.isUserPartOfCollection(ctx, userID, collectionID) {
		return nil, domain.ErrUnauthorized
	}

	secret.UpdatedBy = userID
	secret.UpdatedAt = time.Now()

	updatedSecretDAO, err := svc.secretStorage.UpdateSecret(ctx, converter.ToSecretDAO(secret))
	if err != nil {
		svc.log.Error("Error updating secret:", sl.Err(err))
		return nil, domain.ErrNoUpdatedData
	}

	secret.LinkedSecretId = updatedSecretDAO.LinkedSecretId

	switch secret.SecretType {
	case domain.PasswordSecretType:
		if updatedPasswordSecret, err := svc.updatePasswordSecret(ctx, secret, encryptionKey); err != nil {
			return nil, err
		} else {
			updatedSecretDAO.PasswordSecret = *updatedPasswordSecret
		}
	case domain.TextSecretType:
		if updatedTextSecret, err := svc.updateTextSecret(ctx, secret, encryptionKey); err != nil {
			return nil, err
		} else {
			updatedSecretDAO.TextSecret = *updatedTextSecret
		}
	default:
		return nil, domain.ErrInvalidSecretType
	}

	return converter.ToSecret(updatedSecretDAO), nil
}

func (svc *SecretService) updatePasswordSecret(ctx context.Context, secret *domain.Secret, encryptionKey []byte) (*dao.PasswordSecretDAO, error) {
	if secret.PasswordSecret == nil {
		return nil, domain.ErrInvalidSecretType
	}

	encryptedPassword, err := cipherkit.Encrypt([]byte(secret.PasswordSecret.Password), encryptionKey)
	if err != nil {
		svc.log.Error("Error encrypting password secret:", sl.Err(err))
		return nil, domain.ErrInternal
	}

	passwordSecretDAO := converter.ToPasswordSecretDAO(secret.PasswordSecret)
	passwordSecretDAO.ID = secret.LinkedSecretId
	passwordSecretDAO.Password = encryptedPassword

	updatedPasswordSecretDAO, err := svc.secretStorage.UpdatePasswordSecret(ctx, passwordSecretDAO)
	if err != nil {
		svc.log.Error("Error updating password secret:", sl.Err(err))
		return nil, domain.ErrNoUpdatedData
	}

	updatedPasswordSecretDAO.Password = []byte(secret.PasswordSecret.Password)
	return updatedPasswordSecretDAO, nil
}

func (svc *SecretService) updateTextSecret(ctx context.Context, secret *domain.Secret, encryptionKey []byte) (*dao.TextSecretDAO, error) {
	if secret.TextSecret == nil {
		return nil, domain.ErrInvalidSecretType
	}

	encryptedText, err := cipherkit.Encrypt([]byte(secret.TextSecret.Text), encryptionKey)
	if err != nil {
		svc.log.Error("Error encrypting text secret:", sl.Err(err))
		return nil, domain.ErrInternal
	}

	textSecretDAO := converter.ToTextSecretDAO(secret.TextSecret)
	textSecretDAO.ID = secret.LinkedSecretId
	textSecretDAO.Text = encryptedText

	updatedTextSecretDAO, err := svc.secretStorage.UpdateTextSecret(ctx, textSecretDAO)
	if err != nil {
		svc.log.Error("Error updating text secret:", sl.Err(err))
		return nil, domain.ErrNoUpdatedData
	}

	updatedTextSecretDAO.Text = []byte(secret.TextSecret.Text)
	return updatedTextSecretDAO, nil
}

// DeleteSecret deletes a secret
func (svc *SecretService) DeleteSecret(ctx context.Context, userID, collectionID, secretID uuid.UUID) error {
	if !svc.isUserPartOfCollection(ctx, userID, collectionID) {
		return domain.ErrUnauthorized
	}

	err := svc.secretStorage.DeleteSecret(ctx, secretID)
	if err != nil {
		svc.log.Error("Error deleting secret:", "secretID", secretID, sl.Err(err))
		return domain.ErrInternal
	}

	return nil
}

// isUserPartOfCollection checks if the user is part of the given collection
func (svc *SecretService) isUserPartOfCollection(ctx context.Context, userID, collectionID uuid.UUID) bool {
	isPartOfCollection, err := svc.collectionStorage.IsUserPartOfCollection(ctx, userID, collectionID)
	if err != nil {
		svc.log.Error("Error checking user collection membership:", "userID", userID, "collectionID", collectionID, sl.Err(err))
		return false
	}

	return isPartOfCollection
}
