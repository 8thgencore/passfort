package auth

import (
	"log/slog"
	"time"

	mailGrpc "github.com/8thgencore/passfort/internal/clients/mail/grpc"
	"github.com/8thgencore/passfort/internal/service"
	"github.com/8thgencore/passfort/internal/service/adapters/cache"
	"github.com/8thgencore/passfort/internal/service/adapters/storage"
)

/**
 * AuthService implements service.AuthService interface
 * and provides an access to the user repository
 * and token service
 */
type AuthService struct {
	log             *slog.Logger
	refreshTokenTTL time.Duration
	storage         storage.UserRepository
	cache           cache.CacheRepository
	tokenManager    service.TokenService
	otp             service.OtpService
	mailClient      *mailGrpc.Client
}

// NewAuthService creates a new auth service instance
func NewAuthService(
	log *slog.Logger,
	refreshTokenTTL time.Duration,
	storage storage.UserRepository,
	cache cache.CacheRepository,
	tokenManager service.TokenService,
	otpService service.OtpService,
	mailClient *mailGrpc.Client,

) *AuthService {
	return &AuthService{
		log,
		refreshTokenTTL,
		storage,
		cache,
		tokenManager,
		otpService,
		mailClient,
	}
}
