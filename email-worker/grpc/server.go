package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	grpcserver "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"booking-system/email-worker/config"
	"booking-system/email-worker/models"
	"booking-system/email-worker/processor"
	"booking-system/email-worker/protos"
	"booking-system/email-worker/services"
)

// Server represents the gRPC server
type Server struct {
	protos.UnimplementedEmailServiceServer
	processor    *processor.Processor
	emailService *services.EmailService
	logger       *zap.Logger
	config       *config.Config
	grpcServer   *grpcserver.Server
}

// NewServer creates a new gRPC server
func NewServer(processor *processor.Processor, emailService *services.EmailService, config *config.Config, logger *zap.Logger) *Server {
	return &Server{
		processor:    processor,
		emailService: emailService,
		logger:       logger,
		config:       config,
	}
}

// Start starts the gRPC server
func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.grpcServer = grpcserver.NewServer()
	protos.RegisterEmailServiceServer(s.grpcServer, s)

	// Enable reflection for debugging
	reflection.Register(s.grpcServer)

	s.logger.Info("Starting gRPC server", zap.Int("port", port))

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	return nil
}

// Stop stops the gRPC server
func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.logger.Info("Stopping gRPC server")
		s.grpcServer.GracefulStop()
	}
}

// CreateEmailJob implements the CreateEmailJob gRPC method
func (s *Server) CreateEmailJob(ctx context.Context, req *protos.CreateEmailJobRequest) (*protos.CreateEmailJobResponse, error) {
	s.logger.Info("Creating email job",
		zap.String("template_name", req.TemplateName),
		zap.Strings("recipients", req.To),
	)

	// Create email job
	job := s.createEmailJobFromRequest(req)

	// Publish to queue
	err := s.processor.PublishJob(ctx, job)
	if err != nil {
		s.logger.Error("Failed to publish email job", zap.Error(err))
		return &protos.CreateEmailJobResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create email job: %v", err),
		}, nil
	}

	return &protos.CreateEmailJobResponse{
		Job: &protos.EmailJob{
			Id: job.ID.String(),
		},
		Success: true,
		Message: "Email job created successfully",
	}, nil
}

// GetEmailJob implements the GetEmailJob gRPC method
func (s *Server) GetEmailJob(ctx context.Context, req *protos.GetEmailJobRequest) (*protos.GetEmailJobResponse, error) {
	// This would need to be implemented to query the database
	// For now, return a placeholder response
	return &protos.GetEmailJobResponse{
		Success: true,
		Message: "Email job retrieved successfully",
		Job: &protos.EmailJob{
			Id:     fmt.Sprintf("%d", req.JobId),
			Status: protos.JobStatus_STATUS_PENDING,
		},
	}, nil
}

// UpdateEmailJobStatus implements the UpdateEmailJobStatus gRPC method
func (s *Server) UpdateEmailJobStatus(ctx context.Context, req *protos.UpdateEmailJobStatusRequest) (*protos.UpdateEmailJobStatusResponse, error) {
	// This would need to be implemented to update job status
	return &protos.UpdateEmailJobStatusResponse{
		Success: true,
		Message: "Email job status updated successfully",
	}, nil
}

// ListEmailJobs implements the ListEmailJobs gRPC method
func (s *Server) ListEmailJobs(ctx context.Context, req *protos.ListEmailJobsRequest) (*protos.ListEmailJobsResponse, error) {
	// This would need to be implemented to list jobs
	return &protos.ListEmailJobsResponse{
		Success: true,
		Message: "Email jobs listed successfully",
		Jobs:    []*protos.EmailJob{},
		Total:   0,
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}

// GetEmailTemplate implements the GetEmailTemplate gRPC method
func (s *Server) GetEmailTemplate(ctx context.Context, req *protos.GetEmailTemplateRequest) (*protos.GetEmailTemplateResponse, error) {
	// This would need to be implemented to get templates
	return &protos.GetEmailTemplateResponse{
		Success: true,
		Message: "Email template retrieved successfully",
		Template: &protos.EmailTemplate{
			Id:   req.TemplateId,
			Name: req.Name,
		},
	}, nil
}

// ListEmailTemplates implements the ListEmailTemplates gRPC method
func (s *Server) ListEmailTemplates(ctx context.Context, req *protos.ListEmailTemplatesRequest) (*protos.ListEmailTemplatesResponse, error) {
	// This would need to be implemented to list templates
	return &protos.ListEmailTemplatesResponse{
		Success:   true,
		Message:   "Email templates listed successfully",
		Templates: []*protos.EmailTemplate{},
		Total:     0,
	}, nil
}

// CreateEmailTemplate implements the CreateEmailTemplate gRPC method
func (s *Server) CreateEmailTemplate(ctx context.Context, req *protos.CreateEmailTemplateRequest) (*protos.CreateEmailTemplateResponse, error) {
	// This would need to be implemented to create templates
	return &protos.CreateEmailTemplateResponse{
		Success: true,
		Message: "Email template created successfully",
		Template: &protos.EmailTemplate{
			Name: req.Name,
		},
	}, nil
}

// UpdateEmailTemplate implements the UpdateEmailTemplate gRPC method
func (s *Server) UpdateEmailTemplate(ctx context.Context, req *protos.UpdateEmailTemplateRequest) (*protos.UpdateEmailTemplateResponse, error) {
	// This would need to be implemented to update templates
	return &protos.UpdateEmailTemplateResponse{
		Success: true,
		Message: "Email template updated successfully",
	}, nil
}

// GetEmailTracking implements the GetEmailTracking gRPC method
func (s *Server) GetEmailTracking(ctx context.Context, req *protos.GetEmailTrackingRequest) (*protos.GetEmailTrackingResponse, error) {
	// This would need to be implemented to get tracking info
	return &protos.GetEmailTrackingResponse{
		Success: true,
		Message: "Email tracking retrieved successfully",
	}, nil
}

// UpdateEmailTracking implements the UpdateEmailTracking gRPC method
func (s *Server) UpdateEmailTracking(ctx context.Context, req *protos.UpdateEmailTrackingRequest) (*protos.UpdateEmailTrackingResponse, error) {
	// This would need to be implemented to update tracking info
	return &protos.UpdateEmailTrackingResponse{
		Success: true,
		Message: "Email tracking updated successfully",
	}, nil
}

// SendVerificationEmail implements the SendVerificationEmail gRPC method
func (s *Server) SendVerificationEmail(ctx context.Context, req *protos.SendVerificationEmailRequest) (*protos.SendVerificationEmailResponse, error) {
	s.logger.Info("Sending verification email",
		zap.String("user_id", req.UserId),
		zap.String("user_email", req.UserEmail),
	)

	// Create email job for verification
	job := models.NewEmailJob(
		[]string{req.UserEmail}, // to
		nil,                     // cc
		nil,                     // bcc
		"email_verification",    // templateName
		map[string]any{
			"Name":            req.UserName,
			"PinCode":         req.PinCode,
			"VerificationURL": req.VerificationUrl,
			"ExpiryMinutes":   req.ExpiryTime,
		},
		models.JobPriorityNormal,
	)
	job.SetMaxRetries(3)

	// Save job to database FIRST for tracking
	if err := s.emailService.CreateTrackedEmailJob(ctx, job); err != nil {
		s.logger.Error("Failed to save verification email job to database", zap.Error(err))
		return &protos.SendVerificationEmailResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to save verification email job: %v", err),
		}, nil
	}

	// Then publish to queue
	err := s.processor.PublishJob(ctx, job)
	if err != nil {
		s.logger.Error("Failed to publish verification email job", zap.Error(err))
		return &protos.SendVerificationEmailResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to send verification email: %v", err),
		}, nil
	}

	return &protos.SendVerificationEmailResponse{
		Success: true,
		Message: "Verification email sent successfully",
		JobId:   job.ID.String(),
		PinCode: req.PinCode,
	}, nil
}

// SendVerificationReminder implements the SendVerificationReminder gRPC method
func (s *Server) SendVerificationReminder(ctx context.Context, req *protos.SendVerificationReminderRequest) (*protos.SendVerificationReminderResponse, error) {
	s.logger.Info("Sending verification reminder",
		zap.String("user_id", req.UserId),
		zap.String("user_email", req.UserEmail),
	)

	// Create email job for verification reminder
	job := models.NewEmailJob(
		[]string{req.UserEmail},       // to
		nil,                           // cc
		nil,                           // bcc
		"email_verification_reminder", // templateName
		map[string]any{
			"user_name":        req.UserName,
			"verification_url": req.VerificationUrl,
		},
		models.JobPriorityNormal,
	)
	job.SetMaxRetries(3)

	// Publish to queue
	err := s.processor.PublishJob(ctx, job)
	if err != nil {
		s.logger.Error("Failed to publish verification reminder job", zap.Error(err))
		return &protos.SendVerificationReminderResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to send verification reminder: %v", err),
		}, nil
	}

	return &protos.SendVerificationReminderResponse{
		Success: true,
		Message: "Verification reminder sent successfully",
		JobId:   job.ID.String(),
	}, nil
}

// ResendVerificationEmail implements the ResendVerificationEmail gRPC method
func (s *Server) ResendVerificationEmail(ctx context.Context, req *protos.ResendVerificationEmailRequest) (*protos.SendVerificationEmailResponse, error) {
	s.logger.Info("Resending verification email",
		zap.String("user_id", req.UserId),
		zap.String("user_email", req.UserEmail),
	)

	// Create email job for verification resend
	job := models.NewEmailJob(
		[]string{req.UserEmail}, // to
		nil,                     // cc
		nil,                     // bcc
		"email_verification",    // templateName
		map[string]any{
			"user_name":        req.UserName,
			"verification_url": req.VerificationUrl,
			"is_resend":        true,
		},
		models.JobPriorityNormal,
	)
	job.SetMaxRetries(3)

	// Publish to queue
	err := s.processor.PublishJob(ctx, job)
	if err != nil {
		s.logger.Error("Failed to publish verification resend job", zap.Error(err))
		return &protos.SendVerificationEmailResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to resend verification email: %v", err),
		}, nil
	}

	return &protos.SendVerificationEmailResponse{
		Success: true,
		Message: "Verification email resent successfully",
		JobId:   job.ID.String(),
	}, nil
}

// Health implements the Health gRPC method
func (s *Server) Health(ctx context.Context, req *protos.HealthRequest) (*protos.HealthResponse, error) {
	// Check processor health
	err := s.processor.Health(ctx)
	if err != nil {
		return &protos.HealthResponse{
			Status:    "unhealthy",
			Version:   "1.0.0",
			Timestamp: time.Now().Format(time.RFC3339),
		}, nil
	}

	return &protos.HealthResponse{
		Status:    "healthy",
		Version:   "1.0.0",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// createEmailJobFromRequest creates an EmailJob from a gRPC request
func (s *Server) createEmailJobFromRequest(req *protos.CreateEmailJobRequest) *models.EmailJob {
	// Convert protobuf priority to model priority
	var priority models.JobPriority
	switch req.Priority {
	case protos.JobPriority_PRIORITY_HIGH:
		priority = models.JobPriorityHigh
	case protos.JobPriority_PRIORITY_LOW:
		priority = models.JobPriorityLow
	default:
		priority = models.JobPriorityNormal
	}

	// Convert variables from map[string]string to map[string]any
	variables := make(map[string]any)
	for k, v := range req.Variables {
		variables[k] = v
	}

	// Create email job
	job := models.NewEmailJob(
		req.To,
		req.Cc,
		req.Bcc,
		req.TemplateName,
		variables,
		priority,
	)

	// Set max retries if provided
	if req.MaxRetries > 0 {
		job.MaxRetries = int(req.MaxRetries)
	}

	return job
}
