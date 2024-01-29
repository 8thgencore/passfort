package auth

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/pkg/util"
)

// Login gives a registered user an access token if the credentials are valid
func (as *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := as.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return "", domain.ErrInvalidCredentials
		}
		return "", domain.ErrInternal
	}

	err = util.ComparePassword(password, user.Password)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	accessToken, err := as.ts.CreateToken(user)
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	return accessToken, nil
}
