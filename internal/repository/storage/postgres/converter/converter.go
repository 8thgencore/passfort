package converter

import (
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/dao"
)

func ToUserDAO(user *domain.User) *dao.UserDAO {
	return &dao.UserDAO{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Password:       user.Password,
		MasterPassword: postgres.NullString(user.MasterPassword),
		Salt:           user.Salt,
		IsVerified:     user.IsVerified,
		Role:           string(user.Role),
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func ToUser(userDAO *dao.UserDAO) *domain.User {
	return &domain.User{
		ID:             userDAO.ID,
		Name:           userDAO.Name,
		Email:          userDAO.Email,
		Password:       userDAO.Password,
		MasterPassword: userDAO.MasterPassword.String,
		Salt:           userDAO.Salt,
		IsVerified:     userDAO.IsVerified,
		Role:           domain.UserRoleEnum(userDAO.Role),
		CreatedAt:      userDAO.CreatedAt,
		UpdatedAt:      userDAO.UpdatedAt,
	}
}

func ToCollectionDAO(collection *domain.Collection) *dao.CollectionDAO {
	return &dao.CollectionDAO{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
		CreatedBy:   collection.CreatedBy,
		UpdatedBy:   collection.UpdatedBy,
		CreatedAt:   collection.CreatedAt,
		UpdatedAt:   collection.UpdatedAt,
	}
}

func ToCollection(collectionDAO *dao.CollectionDAO) *domain.Collection {
	return &domain.Collection{
		ID:          collectionDAO.ID,
		Name:        collectionDAO.Name,
		Description: collectionDAO.Description,
		CreatedBy:   collectionDAO.CreatedBy,
		UpdatedBy:   collectionDAO.UpdatedBy,
		CreatedAt:   collectionDAO.CreatedAt,
		UpdatedAt:   collectionDAO.UpdatedAt,
	}
}

// ToSecretDAO converts a domain.Secret to a dao.SecretDAO
func ToSecretDAO(secret *domain.Secret) *dao.SecretDAO {
	secretDAO := &dao.SecretDAO{
		ID:             secret.ID,
		CollectionID:   secret.CollectionID,
		SecretType:     dao.SecretType(secret.SecretType),
		Name:           secret.Name,
		Description:    secret.Description,
		CreatedBy:      secret.CreatedBy,
		UpdatedBy:      secret.UpdatedBy,
		CreatedAt:      secret.CreatedAt,
		UpdatedAt:      secret.UpdatedAt,
		LinkedSecretId: secret.LinkedSecretId,
	}

	switch secret.SecretType {
	case domain.PasswordSecretType:
		if secret.PasswordSecret != nil {
			secretDAO.PasswordSecret = *ToPasswordSecretDAO(secret.PasswordSecret)
		}
	case domain.TextSecretType:
		if secret.TextSecret != nil {
			secretDAO.TextSecret = *ToTextSecretDAO(secret.TextSecret)
		}
	}

	return secretDAO
}

// ToSecret converts a dao.SecretDAO to a domain.Secret
func ToSecret(secretDAO *dao.SecretDAO) *domain.Secret {
	secret := &domain.Secret{
		ID:             secretDAO.ID,
		CollectionID:   secretDAO.CollectionID,
		SecretType:     domain.SecretTypeEnum(secretDAO.SecretType),
		Name:           secretDAO.Name,
		Description:    secretDAO.Description,
		CreatedBy:      secretDAO.CreatedBy,
		UpdatedBy:      secretDAO.UpdatedBy,
		CreatedAt:      secretDAO.CreatedAt,
		UpdatedAt:      secretDAO.UpdatedAt,
		LinkedSecretId: secretDAO.LinkedSecretId,
	}

	switch secretDAO.SecretType {
	case dao.PasswordSecretType:
		secret.PasswordSecret = ToPasswordSecret(&secretDAO.PasswordSecret)
	case dao.TextSecretType:
		secret.TextSecret = ToTextSecret(&secretDAO.TextSecret)
	}

	return secret
}

// ToPasswordSecretDAO converts a domain.PasswordSecret to a dao.PasswordSecretDAO
func ToPasswordSecretDAO(ps *domain.PasswordSecret) *dao.PasswordSecretDAO {
	return &dao.PasswordSecretDAO{
		ID:       ps.ID,
		URL:      ps.URL,
		Login:    ps.Login,
		Password: []byte(ps.Password),
	}
}

// ToPasswordSecret converts a dao.PasswordSecretDAO to a domain.PasswordSecret
func ToPasswordSecret(dao *dao.PasswordSecretDAO) *domain.PasswordSecret {
	return &domain.PasswordSecret{
		ID:       dao.ID,
		URL:      dao.URL,
		Login:    dao.Login,
		Password: string(dao.Password),
	}
}

// ToTextSecretDAO converts a domain.TextSecret to a dao.TextSecretDAO
func ToTextSecretDAO(ts *domain.TextSecret) *dao.TextSecretDAO {
	return &dao.TextSecretDAO{
		ID:   ts.ID,
		Text: []byte(ts.Text),
	}
}

// ToTextSecret converts a dao.TextSecretDAO to a domain.TextSecret
func ToTextSecret(dao *dao.TextSecretDAO) *domain.TextSecret {
	return &domain.TextSecret{
		ID:   dao.ID,
		Text: string(dao.Text),
	}
}
