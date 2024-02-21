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
	secret.CreatedAt = time.Now()
	secret.UpdatedAt = time.Now()

	if !ss.isUserPartOfCollection(ctx, userID, secret.CollectionID) {
		return nil, domain.ErrUnauthorized
	}

	createdSecretDAO, err := ss.secretStorage.CreateSecret(ctx, secret.CollectionID, converter.ToSecretDAO(secret))
	if err != nil {
		ss.log.Error("Error creating secret:", "error", err.Error())
		return nil, err
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
		return nil, err
	}

	var secrets []domain.Secret
	for _, secretDAO := range secretsDAO {
		secrets = append(secrets, *converter.ToSecret(&secretDAO))
	}

	return secrets, nil
}

// GetSecret gets a secret by ID
func (ss *SecretService) GetSecret(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) (*domain.Secret, error) {
	// Implement your business logic and call to the repository's GetSecretByID method
	secretDAO, err := ss.secretStorage.GetSecretByID(ctx, secretID)
	if err != nil {
		ss.log.Error(fmt.Sprintf("Error getting secret %d:", secretID), "error", err.Error())
		return nil, err
	}

	// Check if the user is part of the collection associated with the secret
	// TODO:
	// if !ss.isUserPartOfCollection(ctx, userID, secretDAO.CollectionID) {
	// 	return nil, domain.ErrUnauthorized
	// }

	return converter.ToSecret(secretDAO), nil
}

// UpdateSecret updates a secret
func (ss *SecretService) UpdateSecret(ctx context.Context, userID uuid.UUID, secret *domain.Secret) (*domain.Secret, error) {
	// Implement your business logic, validation, and call to the repository's UpdateSecret method
	// ...

	return nil, nil
}

// DeleteSecret deletes a secret
func (ss *SecretService) DeleteSecret(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) error {
	// Implement your business logic and call to the repository's DeleteSecret method
	// ...

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
