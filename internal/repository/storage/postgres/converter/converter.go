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

func ToSecretDAO(secret *domain.Secret) *dao.SecretDAO {
	return &dao.SecretDAO{
		ID:           secret.ID,
		CollectionID: secret.CollectionID,
		SecretType:   string(secret.SecretType),
		Name:         secret.Name,
		Description:  secret.Description,
		CreatedBy:    secret.CreatedBy,
		UpdatedBy:    secret.UpdatedBy,
		CreatedAt:    secret.CreatedAt,
		UpdatedAt:    secret.UpdatedAt,
	}
}

func ToSecret(secretDAO *dao.SecretDAO) *domain.Secret {
	return &domain.Secret{
		ID:           secretDAO.ID,
		CollectionID: secretDAO.CollectionID,
		SecretType:   domain.SecretTypeEnum(secretDAO.SecretType),
		Name:         secretDAO.Name,
		Description:  secretDAO.Description,
		CreatedBy:    secretDAO.CreatedBy,
		UpdatedBy:    secretDAO.UpdatedBy,
		CreatedAt:    secretDAO.CreatedAt,
		UpdatedAt:    secretDAO.UpdatedAt,
	}
}
