package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"realtime-service/config"
	"realtime-service/pkg/logger"
)

// AuthMiddleware handles JWT authentication for WebSocket connections
type AuthMiddleware struct {
	config *config.Config
}

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: cfg,
	}
}

// ValidateRequest extracts and validates JWT from HTTP request
// Returns user_id if valid, empty string if anonymous allowed, error if invalid
func (m *AuthMiddleware) ValidateRequest(r *http.Request) (string, error) {
	token := m.extractToken(r)

	// Allow anonymous connections if configured
	if token == "" {
		if m.config.WebSocket.AllowAnonymous {
			return "", nil
		}
		return "", errors.New("authentication required")
	}

	// Validate JWT
	claims, err := m.validateToken(token)
	if err != nil {
		return "", err
	}

	return claims.UserID, nil
}

// extractToken gets the JWT token from request
// Checks: Authorization header, query param "token"
func (m *AuthMiddleware) extractToken(r *http.Request) string {
	// Check Authorization header first
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// Check query parameter (for WebSocket connections)
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	return ""
}

// validateToken validates a JWT token and returns claims
func (m *AuthMiddleware) validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.config.JWT.Secret), nil
	})

	if err != nil {
		logger.Debug("JWT validation failed", zap.Error(err))
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	// Check expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	// Validate required claims
	if claims.UserID == "" {
		return nil, errors.New("missing user_id in token")
	}

	return claims, nil
}

// HTTPMiddleware wraps an HTTP handler with authentication
func (m *AuthMiddleware) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := m.ValidateRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user ID to request context if authenticated
		if userID != "" {
			r.Header.Set("X-User-ID", userID)
		}

		next.ServeHTTP(w, r)
	})
}
