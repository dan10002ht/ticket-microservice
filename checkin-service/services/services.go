package services

import (
	"go.uber.org/zap"

	"checkin-service/grpcclient"
	"checkin-service/repositories"
)

type Services struct {
	Checkin *CheckinService
}

func NewServices(repo *repositories.CheckinRepository, clients *grpcclient.Clients, logger *zap.Logger) *Services {
	return &Services{
		Checkin: NewCheckinService(repo, clients, logger),
	}
}

func (s *Services) Close() {
	// nothing to close on the service layer; connection cleanup is in grpcclient
}
