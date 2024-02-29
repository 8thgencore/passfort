package user

import (
	"log/slog"

	mailGrpc "github.com/8thgencore/passfort/internal/clients/mail/grpc"
	"github.com/8thgencore/passfort/internal/service"
	"github.com/8thgencore/passfort/internal/service/adapters/cache"
	"github.com/8thgencore/passfort/internal/service/adapters/storage"
)

/**
 * UserService implements service.UserService interface
 * and provides an access to the user repository
 * and cache service
 */
type UserService struct {
	log        *slog.Logger
	storage    storage.UserRepository
	cache      cache.CacheRepository
	otp        service.OtpService
	mailClient mailGrpc.Client
}

// NewUserService creates a new user service instance
func NewUserService(log *slog.Logger,
	storage storage.UserRepository,
	cache cache.CacheRepository,
	otpService service.OtpService,
	mailClient mailGrpc.Client,
) *UserService {
	return &UserService{
		log,
		storage,
		cache,
		otpService,
		mailClient,
	}
}
