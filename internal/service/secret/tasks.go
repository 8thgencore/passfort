package secret

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/dao"
	"github.com/8thgencore/passfort/pkg/cipherkit"
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
		svc.log.Error("error marshaling payload:", "error", err.Error())
		return domain.ErrInternal
	}

	task := asynq.NewTask(TypeReencryptSecrets, payload)
	if _, err := svc.asynqClient.Enqueue(task); err != nil {
		svc.log.Error("error enqueueing reencrypt secrets task:", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// HandleReencryptSecretsTask handles the re-encryption of all user's secrets
func (svc *SecretService) HandleReencryptSecretsTask(ctx context.Context, t *asynq.Task) error {
	var p ReencryptSecretsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		svc.log.Error("error unmarshaling payload:", "error", err.Error())
		return fmt.Errorf("unmarshal payload: %v", err)
	}

	if !svc.isUserPartOfCollection(ctx, p.UserID, uuid.Nil) {
		return domain.ErrUnauthorized
	}

	// Fetch the collections the user is part of
	collections, err := svc.collectionStorage.ListCollectionsByUserID(ctx, p.UserID, 0, 10)
	if err != nil {
		svc.log.Error("Error fetching collections for user:", "userID", p.UserID, "error", err.Error())
		return domain.ErrDataNotFound
	}

	var secretsDao []dao.SecretDAO
	for _, collection := range collections {
		partSecretsDAO, err := svc.secretStorage.ListSecretsByCollectionID(ctx, collection.ID, 0, 10)
		if err != nil {
			svc.log.Error("Error listing secrets for collection", "error", err.Error(), "collection_id", collection.ID)
			return domain.ErrDataNotFound
		}

		secretsDao = append(secretsDao, partSecretsDAO...)
	}

	for _, secretDAO := range secretsDao {

		// Decrypt and re-encrypt the linked secret based on its type
		switch secretDAO.SecretType {
		case dao.PasswordSecretType:
			decryptPassword, err := cipherkit.Decrypt(secretDAO.PasswordSecret.Password, p.OldEncryptionKey)
			if err != nil {
				svc.log.Error("Error decrypting password secret:", "error", err.Error())
				return domain.ErrInternal
			}

			encryptedPassword, err := cipherkit.Encrypt(decryptPassword, p.NewEncryptionKey)
			if err != nil {
				svc.log.Error("Error encrypting password secret:", "error", err.Error())
				return domain.ErrInternal
			}

			secretDAO.PasswordSecret.Password = encryptedPassword
		case dao.TextSecretType:
			decryptText, err := cipherkit.Decrypt(secretDAO.TextSecret.Text, p.OldEncryptionKey)
			if err != nil {
				svc.log.Error("Error decrypting text secret:", "error", err.Error())
				return domain.ErrInternal
			}

			encryptedText, err := cipherkit.Encrypt(decryptText, p.NewEncryptionKey)
			if err != nil {
				svc.log.Error("Error encrypting text secret:", "error", err.Error())
				return domain.ErrInternal
			}

			secretDAO.TextSecret.Text = encryptedText
		default:
			return domain.ErrInvalidSecretType
		}

		// Update the secret with the new encryption
		if _, err := svc.UpdateSecret(ctx, p.UserID, secretDAO.CollectionID, converter.ToSecret(&secretDAO), p.NewEncryptionKey); err != nil {
			svc.log.Error("Error updating secret:", "error", err.Error())
			return err
		}
	}

	return nil
}
