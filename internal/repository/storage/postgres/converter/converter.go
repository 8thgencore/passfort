package converter

import (
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/dao"
)

func ToUserDAO(user *domain.User) *dao.UserDAO {
	return &dao.UserDAO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUser(userDAO *dao.UserDAO) *domain.User {
	return &domain.User{
		ID:        userDAO.ID,
		Name:      userDAO.Name,
		Email:     userDAO.Email,
		Password:  userDAO.Password,
		Role:      domain.UserRoleEnum(userDAO.Role),
		CreatedAt: userDAO.CreatedAt,
		UpdatedAt: userDAO.UpdatedAt,
	}
}
