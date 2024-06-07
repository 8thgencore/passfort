package secret

import (
	"context"
	"fmt"
	"time"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/dao"
	"github.com/8thgencore/passfort/pkg/cipherkit"
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

	switch secret.SecretType {
	case domain.PasswordSecretType:
		encryptPassword, err := cipherkit.Encrypt([]byte(secret.PasswordSecret.Password), encryptionKey)
		if err != nil {
			svc.log.Error("Error encrypt password secret:", "error", err.Error())
			return nil, domain.ErrInternal
		}

		passwordSecretDAO := converter.ToPasswordSecretDAO(secret.PasswordSecret)
		passwordSecretDAO.Password = encryptPassword

		newSecret, err := svc.secretStorage.CreatePasswordSecret(ctx, passwordSecretDAO)
		if err != nil {
			svc.log.Error("Error creating text secret:", "error", err.Error())
			return nil, domain.ErrInternal
		}
		secretDAO.LinkedSecretId = newSecret.ID
	case domain.TextSecretType:
		encryptText, err := cipherkit.Encrypt([]byte(secret.TextSecret.Text), encryptionKey)
		if err != nil {
			svc.log.Error("Error encrypt text secret:", "error", err.Error())
			return nil, domain.ErrInternal
		}

		textSecretDAO := converter.ToTextSecretDAO(secret.TextSecret)
		textSecretDAO.Text = encryptText

		newSecret, err := svc.secretStorage.CreateTextSecret(ctx, textSecretDAO)
		if err != nil {
			svc.log.Error("Error creating text secret:", "error", err.Error())
			return nil, domain.ErrInternal
		}
		secretDAO.LinkedSecretId = newSecret.ID
	default:
		return nil, domain.ErrInvalidSecretType
	}

	createdSecretDAO, err := svc.secretStorage.CreateSecret(ctx, secret.CollectionID, secretDAO)

	if err != nil {
		svc.log.Error("Error creating secret:", "error", err.Error())
		return nil, domain.ErrInternal
	}

	return converter.ToSecret(createdSecretDAO), nil
}

// ListSecretsByCollectionID lists secrets for a specific collection ID
func (svc *SecretService) ListSecretsByCollectionID(ctx context.Context, userID uuid.UUID, collectionID uuid.UUID, skip, limit uint64) ([]domain.Secret, error) {
	if !svc.isUserPartOfCollection(ctx, userID, collectionID) {
		return nil, domain.ErrUnauthorized
	}

	secretsDAO, err := svc.secretStorage.ListSecretsByCollectionID(ctx, collectionID, skip, limit)
	if err != nil {
		svc.log.Error(fmt.Sprintf("Error listing secrets for collection %d:", collectionID), "error", err.Error())
		return nil, domain.ErrDataNotFound
	}

	var secrets []domain.Secret
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
		secretDAO.LinkedSecret = passwordSecretDAO
	case dao.TextSecretType:
		textSecretDAO, err := svc.getAndDecryptTextSecret(ctx, secretDAO.LinkedSecretId, encryptionKey)
		if err != nil {
			return nil, err
		}
		secretDAO.LinkedSecret = textSecretDAO
	default:
		svc.log.Error("Invalid secret type", "secretType", secretDAO.SecretType, "secretID", secretID)
		return nil, domain.ErrInvalidSecretType
	}

	secret := converter.ToSecret(secretDAO)
	return secret, nil
}

func (svc *SecretService) getAndDecryptPasswordSecret(ctx context.Context, linkedSecretID uuid.UUID, encryptionKey []byte) (*dao.PasswordSecretDAO, error) {
	passwordSecretDAO, err := svc.secretStorage.GetPasswordSecretByID(ctx, linkedSecretID)
	if err != nil {
		svc.log.Error("Error getting password secret by ID", "linkedSecretID", linkedSecretID, "error", err)
		return nil, domain.ErrDataNotFound
	}

	decryptPassword, err := cipherkit.Decrypt(passwordSecretDAO.Password, encryptionKey)
	if err != nil {
		svc.log.Error("Error decrypting password secret", "linkedSecretID", linkedSecretID, "error", err)
		return nil, domain.ErrInternal
	}

	passwordSecretDAO.Password = decryptPassword
	return passwordSecretDAO, nil
}

func (svc *SecretService) getAndDecryptTextSecret(ctx context.Context, linkedSecretID uuid.UUID, encryptionKey []byte) (*dao.TextSecretDAO, error) {
	textSecretDAO, err := svc.secretStorage.GetTextSecretByID(ctx, linkedSecretID)
	if err != nil {
		svc.log.Error("Error getting text secret by ID", "linkedSecretID", linkedSecretID, "error", err)
		return nil, domain.ErrDataNotFound
	}

	decryptText, err := cipherkit.Decrypt(textSecretDAO.Text, encryptionKey)
	if err != nil {
		svc.log.Error("Error decrypting text secret", "linkedSecretID", linkedSecretID, "error", err)
		return nil, domain.ErrInternal
	}

	textSecretDAO.Text = decryptText
	return textSecretDAO, nil
}

// UpdateSecret updates a secret
func (svc *SecretService) UpdateSecret(ctx context.Context, userID, collectionID uuid.UUID, secret *domain.Secret, encryptionKey []byte) (*domain.Secret, error) {
	// Check if the user is part of the collection
	if !svc.isUserPartOfCollection(ctx, userID, collectionID) {
		return nil, domain.ErrUnauthorized
	}

	// Update the fields related to who updated the secret and when
	secret.UpdatedBy = userID
	secret.UpdatedAt = time.Now()

	// Update the main secret
	updatedSecretDAO, err := svc.secretStorage.UpdateSecret(ctx, converter.ToSecretDAO(secret))
	if err != nil {
		svc.log.Error("Error updating secret:", "error", err.Error())
		return nil, domain.ErrNoUpdatedData
	}
	secret.LinkedSecretId = updatedSecretDAO.LinkedSecretId

	// Update linked secret based on secret type
	updatedLinkedSecret, err := svc.updateLinkedSecret(ctx, secret, encryptionKey)
	if err != nil {
		return nil, err
	}
	updatedSecretDAO.LinkedSecret = updatedLinkedSecret

	// Convert the updated dao.SecretDAO back to domain.Secret and return it
	return converter.ToSecret(updatedSecretDAO), nil
}

func (svc *SecretService) updateLinkedSecret(ctx context.Context, secret *domain.Secret, encryptionKey []byte) (dao.ISecret, error) {
	switch secret.SecretType {
	case domain.PasswordSecretType:
		if secret.PasswordSecret != nil {
			encryptedPassword, err := cipherkit.Encrypt([]byte(secret.PasswordSecret.Password), encryptionKey)
			if err != nil {
				svc.log.Error("Error encrypting password secret:", "error", err.Error())
				return nil, domain.ErrInternal
			}

			passwordSecretDAO := converter.ToPasswordSecretDAO(secret.PasswordSecret)
			passwordSecretDAO.ID = secret.LinkedSecretId
			passwordSecretDAO.Password = encryptedPassword

			updatedPasswordSecretDAO, err := svc.secretStorage.UpdatePasswordSecret(ctx, passwordSecretDAO)
			if err != nil {
				svc.log.Error("Error updating password secret:", "error", err.Error())
				return nil, domain.ErrNoUpdatedData
			}

			// Restore plain password for returned value
			updatedPasswordSecretDAO.Password = []byte(secret.PasswordSecret.Password)

			return updatedPasswordSecretDAO, nil
		}
	case domain.TextSecretType:
		if secret.TextSecret != nil {
			encryptedText, err := cipherkit.Encrypt([]byte(secret.TextSecret.Text), encryptionKey)
			if err != nil {
				svc.log.Error("Error encrypting text secret:", "error", err.Error())
				return nil, domain.ErrInternal
			}

			textSecretDAO := converter.ToTextSecretDAO(secret.TextSecret)
			textSecretDAO.ID = secret.LinkedSecretId
			textSecretDAO.Text = encryptedText

			updatedTextSecretDAO, err := svc.secretStorage.UpdateTextSecret(ctx, textSecretDAO)
			if err != nil {
				svc.log.Error("Error updating text secret:", "error", err.Error())
				return nil, domain.ErrNoUpdatedData
			}

			// Restore plain text for returned value
			updatedTextSecretDAO.Text = []byte(secret.TextSecret.Text)

			return updatedTextSecretDAO, nil
		}
	default:
		return nil, domain.ErrInvalidSecretType
	}

	return nil, domain.ErrInvalidSecretType
}

// DeleteSecret deletes a secret
func (svc *SecretService) DeleteSecret(ctx context.Context, userID, collectionID, secretID uuid.UUID) error {
	if !svc.isUserPartOfCollection(ctx, userID, collectionID) {
		return domain.ErrUnauthorized
	}

	err := svc.secretStorage.DeleteSecret(ctx, secretID)
	if err != nil {
		svc.log.Error(fmt.Sprintf("Error deleting secrets %d:", secretID), "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// isUserPartOfCollection checks if the user is part of the given collection
func (svc *SecretService) isUserPartOfCollection(ctx context.Context, userID, collectionID uuid.UUID) bool {
	isPartOfCollection, err := svc.collectionStorage.IsUserPartOfCollection(ctx, userID, collectionID)
	if err != nil {
		svc.log.Error(fmt.Sprintf("Error checking if user %d is part of collection %d:", userID, collectionID), "error", err.Error())
		return false
	}

	return isPartOfCollection
}
