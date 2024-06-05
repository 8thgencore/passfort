package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	mailv1 "github.com/8thgencore/mailfort/protos/gen/go/mail/v1"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api mailv1.MailServiceClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOtps := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOtps...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: mailv1.NewMailServiceClient(cc),
	}, nil
}

func (c *Client) SendConfirmationEmail(ctx context.Context, email, otpCode string) (bool, error) {
	const op = "grpc.SendConfirmationEmail"

	resp, err := c.api.SendConfirmationEmailOTPCode(ctx, &mailv1.SendEmailWithOTPCodeRequest{
		Email:   email,
		OtpCode: otpCode,
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.Success, nil
}

func (c *Client) SendPasswordReset(ctx context.Context, email, otpCode string) (bool, error) {
	const op = "grpc.SendPasswordReset"

	resp, err := c.api.SendPasswordResetOTPCode(ctx, &mailv1.SendEmailWithOTPCodeRequest{
		Email:   email,
		OtpCode: otpCode,
	})
	if err != nil {
		c.log.Error("Failed send reset password code", "error", err.Error())
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.Success, nil
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}
