package grpcclient

import (
	"context"
	"venue-service/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	AuthService    *grpc.ClientConn
	UserService    *grpc.ClientConn
	BookingService *grpc.ClientConn
}

func NewClients(cfg config.GRPCConfig) (*Clients, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	authConn, err := grpc.DialContext(ctx, cfg.AuthServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	userConn, err := grpc.DialContext(ctx, cfg.UserServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	bookingConn, err := grpc.DialContext(ctx, cfg.BookingServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Clients{
		AuthService:    authConn,
		UserService:    userConn,
		BookingService: bookingConn,
	}, nil
}

func (c *Clients) Close() error {
	var errs []error

	if c.AuthService != nil {
		if err := c.AuthService.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if c.UserService != nil {
		if err := c.UserService.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if c.BookingService != nil {
		if err := c.BookingService.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs[0] // Return first error
	}
	return nil
}
