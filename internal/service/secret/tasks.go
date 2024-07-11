package secret

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/pkg/logger/sl"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// Task types
const (
	TypeReencryptSecrets = "reencrypt:secrets"
)

// Payload structure for reencrypt secrets task
type ReencryptSecretsPayload struct {
	UserID           uuid.UUID
	OldEncryptionKey []byte
	NewEncryptionKey []byte
}

// ReencryptAllSecrets initiates the re-encryption of all user's secrets
func (svc *SecretService) ReencryptAllSecrets(ctx context.Context, userID uuid.UUID, oldEncryptionKey, newEncryptionKey []byte) error {
	payload, err := json.Marshal(ReencryptSecretsPayload{
		UserID:           userID,
		OldEncryptionKey: oldEncryptionKey,
		NewEncryptionKey: newEncryptionKey,
	})
	if err != nil {
		svc.log.Error("Error marshaling payload:", sl.Err(err))
		return domain.ErrInternal
	}

	task := asynq.NewTask(TypeReencryptSecrets, payload)
	if _, err := svc.asynqClient.Enqueue(task); err != nil {
		svc.log.Error("Error enqueueing reencrypt secrets task:", sl.Err(err))
		return domain.ErrInternal
	}

	return nil
}

// HandleReencryptSecretsTask handles the re-encryption of all user's secrets
func (svc *SecretService) HandleReencryptSecretsTask(ctx context.Context, t *asynq.Task) error {
	var p ReencryptSecretsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		svc.log.Error("Error unmarshaling payload:", sl.Err(err))
		return fmt.Errorf("unmarshal payload: %v", err)
	}

	collections, err := svc.collectionStorage.ListCollectionsByUserID(ctx, p.UserID, 1, 10)
	if err != nil {
		svc.log.Error("Error fetching collections for user:", "userID", p.UserID, sl.Err(err))
		return domain.ErrDataNotFound
	}

	for _, collection := range collections {
		if err := svc.reencryptCollectionSecrets(ctx, collection.ID, p.UserID, p.OldEncryptionKey, p.NewEncryptionKey); err != nil {
			return err
		}
	}

	return nil
}

func (svc *SecretService) reencryptCollectionSecrets(ctx context.Context, collectionID, userID uuid.UUID, oldEncryptionKey, newEncryptionKey []byte) error {
	secretsDAO, err := svc.secretStorage.ListSecretsByCollectionID(ctx, collectionID, 1, 10)
	if err != nil {
		svc.log.Error("Error listing secrets for collection", "collectionID", collectionID, sl.Err(err))
		return domain.ErrDataNotFound
	}

	for _, secretDAO := range secretsDAO {
		secret, err := svc.GetSecret(ctx, userID, secretDAO.CollectionID, secretDAO.ID, oldEncryptionKey)
		if err != nil {
			svc.log.Error("Error getting secret by ID", "secretID", secretDAO.ID, "error", err)
			return domain.ErrDataNotFound
		}

		if _, err := svc.UpdateSecret(ctx, userID, secretDAO.CollectionID, secret, newEncryptionKey); err != nil {
			svc.log.Error("Error updating secret:", "secretID", secretDAO.ID, sl.Err(err))
			return err
		}
	}

	return nil
}
