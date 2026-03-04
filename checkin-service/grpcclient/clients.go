package grpcclient

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"checkin-service/config"
	ticketpb "checkin-service/internal/protos/ticket"
)

type Clients struct {
	Ticket ticketpb.TicketServiceClient
}

func NewClients(cfg *config.Config) (*Clients, error) {
	addr := fmt.Sprintf("%s:%d", cfg.TicketService.Host, cfg.TicketService.Port)

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to ticket-service at %s: %w", addr, err)
	}

	return &Clients{
		Ticket: ticketpb.NewTicketServiceClient(conn),
	}, nil
}
