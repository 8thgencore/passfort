package auth

import (
	"log/slog"

	mailGrpc "github.com/8thgencore/passfort/internal/clients/mail/grpc"
	"github.com/8thgencore/passfort/internal/service"
	"github.com/8thgencore/passfort/internal/service/adapters/storage"
)

/**
 * AuthService implements service.AuthService interface
 * and provides an access to the user repository
 * and token service
 */
type AuthService struct {
	log        *slog.Logger
	storage    storage.UserRepository
	ts         service.TokenService
	otp        service.OtpService
	mailClient mailGrpc.Client
}

// NewAuthService creates a new auth service instance
func NewAuthService(log *slog.Logger,
	storage storage.UserRepository,
	ts service.TokenService,
	otpService service.OtpService,
	mailClient mailGrpc.Client,

) *AuthService {
	return &AuthService{
		log,
		storage,
		ts,
		otpService,
		mailClient,
	}
}
