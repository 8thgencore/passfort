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
	secret.CreatedBy = userID
	secret.UpdatedBy = userID

	if !ss.isUserPartOfCollection(ctx, userID, secret.CollectionID) {
		return nil, domain.ErrUnauthorized
	}

	createdSecretDAO, err := ss.secretStorage.CreateSecret(ctx, secret.CollectionID, converter.ToSecretDAO(secret))

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

	secretDAO, err := ss.secretStorage.GetSecretByID(ctx, secretID)
	if err != nil {
		ss.log.Error(fmt.Sprintf("Error getting secret %d:", secretID), "error", err.Error())
		return nil, domain.ErrDataNotFound
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
