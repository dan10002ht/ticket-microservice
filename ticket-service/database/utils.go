package database

import (
	"fmt"
	"time"
)

// UpdateUpdatedAtColumn is a PostgreSQL function to update the updated_at timestamp
const UpdateUpdatedAtColumn = `
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';
`

// CreateUpdateTrigger creates a trigger to automatically update the updated_at column
func CreateUpdateTrigger(tableName string) string {
	return `
CREATE TRIGGER update_` + tableName + `_updated_at 
    BEFORE UPDATE ON ` + tableName + ` 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
`
}

// GenerateTicketNumber generates a unique ticket number
func GenerateTicketNumber(eventID, seatID string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("TKT-%s-%s-%d", eventID[:8], seatID[:8], timestamp)
}

// GenerateBookingSessionToken generates a unique booking session token
func GenerateBookingSessionToken(userID string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("BKG-%s-%d", userID[:8], timestamp)
}

// GenerateReservationToken generates a unique reservation token
func GenerateReservationToken(bookingSessionID, seatID string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("RSV-%s-%s-%d", bookingSessionID[:8], seatID[:8], timestamp)
}

// CalculateExpirationTime calculates expiration time for bookings and reservations
func CalculateExpirationTime(timeoutMinutes int) time.Time {
	return time.Now().Add(time.Duration(timeoutMinutes) * time.Minute)
}

// IsExpired checks if a timestamp is expired
func IsExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}

// FormatTime formats time for database storage
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTime parses time from database
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}
