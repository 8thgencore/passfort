package secret

import (
	"context"
	"fmt"
	"time"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/google/uuid"
)

// CreateSecret creates a new secret
func (ss *SecretService) CreateSecret(ctx context.Context, userID uuid.UUID, secret *domain.Secret) (*domain.Secret, error) {
	if !ss.isUserPartOfCollection(ctx, userID, secret.CollectionID) {
		return nil, domain.ErrUnauthorized
	}

	secret.CreatedBy = userID
	secret.UpdatedBy = userID
	secretDAO := converter.ToSecretDAO(secret)

	switch secret.SecretType {
	case domain.PasswordSecretType:
		newSecret, err := ss.secretStorage.CreatePasswordSecret(ctx, converter.ToPasswordSecretDAO(secret.PasswordSecret))
		if err != nil {
			ss.log.Error("Error creating password secret:", "error", err.Error())
			return nil, domain.ErrInternal
		}
		secretDAO.LinkedSecretId = newSecret.ID
	case domain.TextSecretType:
		newSecret, err := ss.secretStorage.CreateTextSecret(ctx, converter.ToTextSecretDAO(secret.TextSecret))
		if err != nil {
			ss.log.Error("Error creating text secret:", "error", err.Error())
			return nil, domain.ErrInternal
		}
		secretDAO.LinkedSecretId = newSecret.ID
	default:
		return nil, domain.ErrInvalidSecretType
	}

	createdSecretDAO, err := ss.secretStorage.CreateSecret(ctx, secret.CollectionID, secretDAO)

	if err != nil {
		ss.log.Error("Error creating secret:", "error", err.Error())
		return nil, domain.ErrInternal
	}

	return converter.ToSecret(createdSecretDAO), nil
}

// ListSecretsByCollectionID lists secrets for a specific collection ID
func (ss *SecretService) ListSecretsByCollectionID(ctx context.Context, userID uuid.UUID, collectionID uuid.UUID, skip, limit uint64) ([]domain.Secret, error) {
	if !ss.isUserPartOfCollection(ctx, userID, collectionID) {
		return nil, domain.ErrUnauthorized
	}

	secretsDAO, err := ss.secretStorage.ListSecretsByCollectionID(ctx, collectionID, skip, limit)
	if err != nil {
		ss.log.Error(fmt.Sprintf("Error listing secrets for collection %d:", collectionID), "error", err.Error())
		return nil, domain.ErrDataNotFound
	}

	var secrets []domain.Secret
	for _, secretDAO := range secretsDAO {
		secrets = append(secrets, *converter.ToSecret(&secretDAO))
	}

	return secrets, nil
}

// GetSecret gets a secret by ID
func (ss *SecretService) GetSecret(ctx context.Context, userID, collectionID, secretID uuid.UUID) (*domain.Secret, error) {
	if !ss.isUserPartOfCollection(ctx, userID, collectionID) {
		return nil, domain.ErrUnauthorized
	}

	// Getting the main secret
	secretDAO, err := ss.secretStorage.GetSecretByID(ctx, secretID)
	if err != nil {
		ss.log.Error(fmt.Sprintf("Error getting secret %d:", secretID), "error", err.Error())
		return nil, domain.ErrDataNotFound
	}

	// Depending on the type of secret, get additional data
	switch secretDAO.SecretType {
	case string(domain.PasswordSecretType):
		passwordSecretDAO, err := ss.secretStorage.GetPasswordSecretByID(ctx, secretDAO.LinkedSecretId)
		if err != nil {
			ss.log.Error(fmt.Sprintf("Error getting password secret %s:", secretDAO.LinkedSecretId), "error", err.Error())
			return nil, domain.ErrDataNotFound
		}
		secretDAO.LinkedSecret = passwordSecretDAO
	case string(domain.TextSecretType):
		textSecretDAO, err := ss.secretStorage.GetTextSecretByID(ctx, secretDAO.LinkedSecretId)
		if err != nil {
			ss.log.Error(fmt.Sprintf("Error getting text secret %s:", secretDAO.LinkedSecretId), "error", err.Error())
			return nil, domain.ErrDataNotFound
		}
		secretDAO.LinkedSecret = textSecretDAO
	default:
		ss.log.Error(fmt.Sprintf("Invalid secret type %s for secret %s:", secretDAO.SecretType, secretID))
		return nil, domain.ErrInvalidSecretType
	}

	return converter.ToSecret(secretDAO), nil
}

// UpdateSecret updates a secret
func (ss *SecretService) UpdateSecret(ctx context.Context, userID, collectionID uuid.UUID, secret *domain.Secret) (*domain.Secret, error) {
	// Check if the user is part of the collection
	if !ss.isUserPartOfCollection(ctx, userID, collectionID) {
		return nil, domain.ErrUnauthorized
	}

	// Update the fields related to who updated the secret and when
	secret.UpdatedBy = userID
	secret.UpdatedAt = time.Now()

	// Convert the domain.Secret to dao.SecretDAO
	secretDAO := converter.ToSecretDAO(secret)

	// Call the repository to update the secret
	updatedSecretDAO, err := ss.secretStorage.UpdateSecret(ctx, secretDAO)
	if err != nil {
		ss.log.Error("Error updating secret:", "error", err.Error())
		return nil, domain.ErrNoUpdatedData
	}

	switch secret.SecretType {
	case domain.PasswordSecretType:
		if secret.PasswordSecret != nil {
			passwordSecretDAO := converter.ToPasswordSecretDAO(&domain.PasswordSecret{
				ID:       updatedSecretDAO.LinkedSecretId,
				URL:      secret.PasswordSecret.URL,
				Login:    secret.PasswordSecret.Login,
				Password: secret.PasswordSecret.Password,
			})
			updatedPasswordSecretDAO, err := ss.secretStorage.UpdatePasswordSecret(ctx, passwordSecretDAO)
			if err != nil {
				ss.log.Error("Error updating password secret:", "error", err.Error())
				return nil, domain.ErrNoUpdatedData
			}
			updatedSecretDAO.LinkedSecret = updatedPasswordSecretDAO
		}
	case domain.TextSecretType:
		if secret.TextSecret != nil {
			textSecretDAO := converter.ToTextSecretDAO(&domain.TextSecret{
				ID:   updatedSecretDAO.LinkedSecretId,
				Text: secret.TextSecret.Text,
			})
			updatedTextSecretDAO, err := ss.secretStorage.UpdateTextSecret(ctx, textSecretDAO)
			if err != nil {
				ss.log.Error("Error updating text secret:", "error", err.Error())
				return nil, domain.ErrNoUpdatedData
			}
			updatedSecretDAO.LinkedSecret = updatedTextSecretDAO
		}
	default:
		return nil, domain.ErrInvalidSecretType
	}

	// Convert the updated dao.SecretDAO back to domain.Secret and return it
	return converter.ToSecret(updatedSecretDAO), nil
}

// DeleteSecret deletes a secret
func (ss *SecretService) DeleteSecret(ctx context.Context, userID, collectionID, secretID uuid.UUID) error {
	if !ss.isUserPartOfCollection(ctx, userID, collectionID) {
		return domain.ErrUnauthorized
	}

	err := ss.secretStorage.DeleteSecret(ctx, secretID)
	if err != nil {
		ss.log.Error(fmt.Sprintf("Error deleting secrets %d:", secretID), "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// isUserPartOfCollection checks if the user is part of the given collection
func (ss *SecretService) isUserPartOfCollection(ctx context.Context, userID, collectionID uuid.UUID) bool {
	isPartOfCollection, err := ss.collectionStorage.IsUserPartOfCollection(ctx, userID, collectionID)
	if err != nil {
		ss.log.Error(fmt.Sprintf("Error checking if user %d is part of collection %d:", userID, collectionID), "error", err.Error())
		return false
	}

	return isPartOfCollection
}
