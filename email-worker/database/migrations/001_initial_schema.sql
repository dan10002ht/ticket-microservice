-- Migration: 001_initial_schema.sql
-- Description: Initial database schema for email worker service with Hybrid ID Pattern
-- Created: 2024-01-01

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Email Jobs Table with Hybrid ID Pattern
CREATE TABLE IF NOT EXISTS email_jobs (
    id BIGSERIAL PRIMARY KEY,                           -- Internal ID for performance
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(), -- Public ID for API
    to_emails TEXT[] NOT NULL, -- Array of recipient emails
    cc_emails TEXT[], -- Array of CC emails
    bcc_emails TEXT[], -- Array of BCC emails
    template_name VARCHAR(255) NOT NULL,
    variables JSONB, -- Template variables
    status VARCHAR(20) DEFAULT 'pending',
    priority INTEGER DEFAULT 2, -- 1=high, 2=normal, 3=low
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    error_message TEXT,
    processed_at TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email Templates Table (using string ID for template names)
CREATE TABLE IF NOT EXISTS email_templates (
    id VARCHAR(100) PRIMARY KEY, -- Template identifier (e.g., 'email_verification', 'booking_confirmation')
    name VARCHAR(255) NOT NULL,
    subject VARCHAR(500),
    html_template TEXT,
    text_template TEXT,
    variables JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email Tracking Table with Hybrid ID Pattern
CREATE TABLE IF NOT EXISTS email_tracking (
    id BIGSERIAL PRIMARY KEY,                           -- Internal ID for performance
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(), -- Public ID for API
    job_id BIGINT NOT NULL REFERENCES email_jobs(id) ON DELETE CASCADE,
    provider VARCHAR(50),
    message_id VARCHAR(255),
    status VARCHAR(50),
    sent_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    opened_at TIMESTAMP WITH TIME ZONE,
    clicked_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    bounce_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX idx_email_jobs_public_id ON email_jobs(public_id);
CREATE INDEX idx_email_jobs_status ON email_jobs(status);
CREATE INDEX idx_email_jobs_priority ON email_jobs(priority);
CREATE INDEX idx_email_jobs_created_at ON email_jobs(created_at);
CREATE INDEX idx_email_jobs_template_name ON email_jobs(template_name);
CREATE INDEX idx_email_jobs_to_emails ON email_jobs USING GIN(to_emails);
CREATE INDEX idx_email_jobs_processed_at ON email_jobs(processed_at);
CREATE INDEX idx_email_jobs_sent_at ON email_jobs(sent_at);

CREATE INDEX idx_email_templates_is_active ON email_templates(is_active);
CREATE INDEX idx_email_templates_name ON email_templates(name);

CREATE INDEX idx_email_tracking_public_id ON email_tracking(public_id);
CREATE INDEX idx_email_tracking_job_id ON email_tracking(job_id);
CREATE INDEX idx_email_tracking_status ON email_tracking(status);
CREATE INDEX idx_email_tracking_provider ON email_tracking(provider);
CREATE INDEX idx_email_tracking_created_at ON email_tracking(created_at);

-- Insert default templates
INSERT INTO email_templates (id, name, subject, html_template, text_template, variables) VALUES
(
    'email_verification',
    'Email Verification',
    'Verify your email address',
    '<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Email Verification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background-color: #ffffff;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            font-size: 24px;
            font-weight: bold;
            color: #2c3e50;
            margin-bottom: 10px;
        }
        .title {
            color: #2c3e50;
            font-size: 20px;
            margin-bottom: 20px;
        }
        .pin-container {
            text-align: center;
            margin: 30px 0;
            padding: 20px;
            background-color: #ecf0f1;
            border-radius: 8px;
        }
        .pin-code {
            font-size: 32px;
            font-weight: bold;
            color: #e74c3c;
            letter-spacing: 5px;
            font-family: "Courier New", monospace;
        }
        .pin-label {
            font-size: 14px;
            color: #7f8c8d;
            margin-top: 10px;
        }
        .info {
            background-color: #e8f4fd;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            border-left: 4px solid #3498db;
        }
        .warning {
            background-color: #fff3cd;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            border-left: 4px solid #ffc107;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #ecf0f1;
            color: #7f8c8d;
            font-size: 12px;
        }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #3498db;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            margin: 10px 0;
        }
        .button:hover {
            background-color: #2980b9;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">Booking System</div>
            <h1 class="title">Email Verification</h1>
        </div>

        <p>Hello <strong>{{.Name}}</strong>,</p>

        <p>Thank you for registering with Booking System. To complete your registration, please verify your email address using the verification code below:</p>

        <div class="pin-container">
            <div class="pin-code">{{.PinCode}}</div>
            <div class="pin-label">Your verification code</div>
        </div>

        <div class="info">
            <strong>How to verify:</strong>
            <ul>
                <li>Enter this code in the verification form on our website</li>
                <li>Or click the verification link below</li>
            </ul>
        </div>

        <div style="text-align: center">
            <a href="{{.VerificationURL}}" class="button">Verify Email Address</a>
        </div>

        <div class="warning">
            <strong>Important:</strong>
            <ul>
                <li>This code will expire in 15 minutes</li>
                <li>If you did not request this verification, please ignore this email</li>
                <li>Never share this code with anyone</li>
            </ul>
        </div>

        <div class="footer">
            <p>This is an automated email from Booking System. Please do not reply to this email.</p>
            <p>If you have any questions, please contact our support team.</p>
        </div>
    </div>
</body>
</html>',
    'Hello {{.Name}},

Thank you for registering with Booking System. To complete your registration, please verify your email address using the verification code below:

Your verification code: {{.PinCode}}

How to verify:
- Enter this code in the verification form on our website
- Or visit: {{.VerificationURL}}

Important:
- This code will expire in 15 minutes
- If you did not request this verification, please ignore this email
- Never share this code with anyone

This is an automated email from Booking System. Please do not reply to this email.

If you have any questions, please contact our support team.',
    '{"Name": "string", "PinCode": "string", "VerificationURL": "string"}'
),
(
    'booking_confirmation',
    'Booking Confirmation',
    'Your booking has been confirmed',
    '<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Booking Confirmation</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background-color: #ffffff;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            font-size: 24px;
            font-weight: bold;
            color: #2c3e50;
            margin-bottom: 10px;
        }
        .title {
            color: #2c3e50;
            font-size: 20px;
            margin-bottom: 20px;
        }
        .booking-details {
            background-color: #f8f9fa;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
        }
        .booking-details h3 {
            margin-top: 0;
            color: #2c3e50;
        }
        .detail-row {
            display: flex;
            justify-content: space-between;
            margin: 10px 0;
            padding: 5px 0;
            border-bottom: 1px solid #e9ecef;
        }
        .detail-label {
            font-weight: bold;
            color: #6c757d;
        }
        .detail-value {
            color: #2c3e50;
        }
        .success {
            background-color: #d4edda;
            color: #155724;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            border-left: 4px solid #28a745;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #ecf0f1;
            color: #7f8c8d;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">Booking System</div>
            <h1 class="title">Booking Confirmation</h1>
        </div>

        <div class="success">
            <strong>✅ Your booking has been confirmed!</strong>
            <p>Thank you for your booking. We have received your payment and your tickets are confirmed.</p>
        </div>

        <div class="booking-details">
            <h3>Booking Details</h3>
            <div class="detail-row">
                <span class="detail-label">Booking ID:</span>
                <span class="detail-value">{{.BookingID}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Event:</span>
                <span class="detail-value">{{.EventName}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Date & Time:</span>
                <span class="detail-value">{{.EventDateTime}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Venue:</span>
                <span class="detail-value">{{.VenueName}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Seats:</span>
                <span class="detail-value">{{.Seats}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Total Amount:</span>
                <span class="detail-value">{{.TotalAmount}}</span>
            </div>
        </div>

        <p><strong>What happens next?</strong></p>
        <ul>
            <li>Your tickets will be sent to this email address</li>
            <li>You can also download your tickets from your account</li>
            <li>Please arrive at least 30 minutes before the event</li>
            <li>Bring a valid ID and your tickets</li>
    </ul>

        <div class="footer">
            <p>This is an automated email from Booking System. Please do not reply to this email.</p>
            <p>If you have any questions, please contact our support team.</p>
        </div>
    </div>
</body>
</html>',
    'Your booking has been confirmed!

Booking Details:
- Booking ID: {{.BookingID}}
- Event: {{.EventName}}
- Date & Time: {{.EventDateTime}}
- Venue: {{.VenueName}}
- Seats: {{.Seats}}
- Total Amount: {{.TotalAmount}}

What happens next?
- Your tickets will be sent to this email address
- You can also download your tickets from your account
- Please arrive at least 30 minutes before the event
- Bring a valid ID and your tickets

This is an automated email from Booking System. Please do not reply to this email.

If you have any questions, please contact our support team.',
    '{"BookingID": "string", "EventName": "string", "EventDateTime": "string", "VenueName": "string", "Seats": "string", "TotalAmount": "string"}'
),
(
    'password_reset',
    'Password Reset Request',
    'Reset your password',
    '<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Password Reset</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background-color: #ffffff;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            font-size: 24px;
            font-weight: bold;
            color: #2c3e50;
            margin-bottom: 10px;
        }
        .title {
            color: #2c3e50;
            font-size: 20px;
            margin-bottom: 20px;
        }
        .warning {
            background-color: #fff3cd;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            border-left: 4px solid #ffc107;
        }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #dc3545;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            margin: 10px 0;
        }
        .button:hover {
            background-color: #c82333;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #ecf0f1;
            color: #7f8c8d;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">Booking System</div>
            <h1 class="title">Password Reset Request</h1>
        </div>

        <p>Hello <strong>{{.Name}}</strong>,</p>

        <p>We received a request to reset your password for your Booking System account. If you made this request, please click the button below to reset your password:</p>

        <div style="text-align: center">
            <a href="{{.ResetURL}}" class="button">Reset Password</a>
        </div>

        <div class="warning">
            <strong>Important:</strong>
            <ul>
                <li>This link will expire in 1 hour</li>
                <li>If you did not request a password reset, please ignore this email</li>
                <li>Never share this link with anyone</li>
            </ul>
        </div>

        <p>If the button above does not work, you can copy and paste this link into your browser:</p>
        <p style="word-break: break-all; color: #007bff;">{{.ResetURL}}</p>

        <div class="footer">
            <p>This is an automated email from Booking System. Please do not reply to this email.</p>
            <p>If you have any questions, please contact our support team.</p>
        </div>
    </div>
</body>
</html>',
    'Hello {{.Name}},

We received a request to reset your password for your Booking System account. If you made this request, please click the link below to reset your password:

{{.ResetURL}}

Important:
- This link will expire in 1 hour
- If you did not request a password reset, please ignore this email
- Never share this link with anyone

This is an automated email from Booking System. Please do not reply to this email.

If you have any questions, please contact our support team.',
    '{"Name": "string", "ResetURL": "string"}'
);

-- Add updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to update updated_at timestamp for email_jobs
CREATE TRIGGER update_email_jobs_updated_at 
    BEFORE UPDATE ON email_jobs 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger to update updated_at timestamp for email_templates
CREATE TRIGGER update_email_templates_updated_at 
    BEFORE UPDATE ON email_templates 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column(); 