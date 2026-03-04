package grpcclient

import (
	"fmt"

	"grpctls"

	"google.golang.org/grpc"

	"checkin-service/config"
	ticketpb "checkin-service/internal/protos/ticket"
)

type Clients struct {
	Ticket ticketpb.TicketServiceClient
}

func NewClients(cfg *config.Config) (*Clients, error) {
	addr := fmt.Sprintf("%s:%d", cfg.TicketService.Host, cfg.TicketService.Port)

	conn, err := grpc.NewClient(addr,
		grpctls.DialOption(),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to ticket-service at %s: %w", addr, err)
	}

	return &Clients{
		Ticket: ticketpb.NewTicketServiceClient(conn),
	}, nil
}
