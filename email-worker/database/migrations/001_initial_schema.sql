-- Migration: 001_initial_schema.sql
-- Description: Initial database schema for email worker service
-- Created: 2024-01-01

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Email Jobs Table
CREATE TABLE IF NOT EXISTS email_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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
    processed_at TIMESTAMP,
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Email Templates Table
CREATE TABLE IF NOT EXISTS email_templates (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    subject VARCHAR(500),
    html_template TEXT,
    text_template TEXT,
    variables JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Email Tracking Table
CREATE TABLE IF NOT EXISTS email_tracking (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id UUID REFERENCES email_jobs(id) ON DELETE CASCADE,
    provider VARCHAR(50),
    message_id VARCHAR(255),
    status VARCHAR(50),
    sent_at TIMESTAMP,
    delivered_at TIMESTAMP,
    opened_at TIMESTAMP,
    clicked_at TIMESTAMP,
    error_message TEXT,
    bounce_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for better performance
CREATE INDEX idx_email_jobs_status ON email_jobs(status);
CREATE INDEX idx_email_jobs_priority ON email_jobs(priority);
CREATE INDEX idx_email_jobs_created_at ON email_jobs(created_at);
CREATE INDEX idx_email_jobs_template_name ON email_jobs(template_name);
CREATE INDEX idx_email_jobs_to_emails ON email_jobs USING GIN(to_emails);
CREATE INDEX idx_email_jobs_processed_at ON email_jobs(processed_at);
CREATE INDEX idx_email_jobs_sent_at ON email_jobs(sent_at);

CREATE INDEX idx_email_templates_is_active ON email_templates(is_active);
CREATE INDEX idx_email_templates_name ON email_templates(name);

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
                <li>This code will expire in <strong>{{.ExpiryMinutes}}</strong> minutes</li>
                <li>If you did not create an account, please ignore this email</li>
                <li>Never share this code with anyone</li>
            </ul>
        </div>

        <p>If you have any questions, please contact our support team.</p>

        <p>Best regards,<br><strong>Booking System Team</strong></p>

        <div class="footer">
            <p>This email was sent for email verification purposes.</p>
            <p>&copy; 2024 Booking System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>',
    'Email Verification - Booking System

Hello {{.Name}},

Thank you for registering with Booking System. To complete your registration, please verify your email address using the verification code below:

Your verification code: {{.PinCode}}

How to verify:
- Enter this code in the verification form on our website
- Or visit: {{.VerificationURL}}

Important:
- This code will expire in {{.ExpiryMinutes}} minutes
- If you did not create an account, please ignore this email
- Never share this code with anyone

If you have any questions, please contact our support team.

Best regards,
Booking System Team

---
This email was sent for email verification purposes.
© 2024 Booking System. All rights reserved.',
    '{"Name": "string", "PinCode": "string", "VerificationURL": "string", "ExpiryMinutes": "number"}'
),
(
    'password_reset',
    'Password Reset',
    'Reset your password',
    '<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Password Reset</title>
</head>
<body>
    <h1>Password Reset Request</h1>
    <p>Hi {{.Email}},</p>
    <p>You requested a password reset. Click the link below to reset your password:</p>
    <a href="{{.ResetURL}}">Reset Password</a>
    <p>This link will expire in 15 minutes.</p>
    <p>If you did not request this, please ignore this email.</p>
    <p>Best regards,<br>Booking System Team</p>
</body>
</html>',
    'Password Reset Request

Hi {{.Email}},

You requested a password reset. Click the link below to reset your password:
{{.ResetURL}}

This link will expire in 15 minutes.

If you did not request this, please ignore this email.

Best regards,
Booking System Team',
    '{"Email": "string", "ResetURL": "string"}'
),
(
    'welcome_email',
    'Welcome Email',
    'Welcome to Booking System!',
    '<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Welcome</title>
</head>
<body>
    <h1>Welcome to Booking System!</h1>
    <p>Hi {{.Name}},</p>
    <p>Thank you for joining Booking System. We are excited to have you on board!</p>
    <p>You can now start booking events and managing your account.</p>
    <p>Best regards,<br>Booking System Team</p>
</body>
</html>',
    'Welcome to Booking System!

Hi {{.Name}},

Thank you for joining Booking System. We are excited to have you on board!

You can now start booking events and managing your account.

Best regards,
Booking System Team',
    '{"Name": "string"}'
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
</head>
<body>
    <h1>Booking Confirmation</h1>
    <p>Hi {{.Name}},</p>
    <p>Your booking has been confirmed!</p>
    <h2>Booking Details:</h2>
    <ul>
        <li><strong>Event:</strong> {{.EventName}}</li>
        <li><strong>Date:</strong> {{.EventDate}}</li>
        <li><strong>Time:</strong> {{.EventTime}}</li>
        <li><strong>Venue:</strong> {{.Venue}}</li>
        <li><strong>Ticket Quantity:</strong> {{.TicketQuantity}}</li>
        <li><strong>Total Amount:</strong> {{.TotalAmount}}</li>
    </ul>
    <p>Booking ID: {{.BookingID}}</p>
    <p>Best regards,<br>Booking System Team</p>
</body>
</html>',
    'Booking Confirmation

Hi {{.Name}},

Your booking has been confirmed!

Booking Details:
- Event: {{.EventName}}
- Date: {{.EventDate}}
- Time: {{.EventTime}}
- Venue: {{.Venue}}
- Ticket Quantity: {{.TicketQuantity}}
- Total Amount: {{.TotalAmount}}

Booking ID: {{.BookingID}}

Best regards,
Booking System Team',
    '{"Name": "string", "EventName": "string", "EventDate": "string", "EventTime": "string", "Venue": "string", "TicketQuantity": "number", "TotalAmount": "string", "BookingID": "string"}'
),
(
    'organization_invitation',
    'Organization Invitation',
    'You have been invited to join an organization',
    '<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Organization Invitation</title>
</head>
<body>
    <h1>Organization Invitation</h1>
    <p>Hi {{.Name}},</p>
    <p>You have been invited to join <strong>{{.OrganizationName}}</strong> on Booking System.</p>
    <p>Role: {{.Role}}</p>
    <p>Click the link below to accept the invitation:</p>
    <a href="{{.InvitationURL}}">Accept Invitation</a>
    <p>This invitation will expire in {{.ExpiryDays}} days.</p>
    <p>Best regards,<br>Booking System Team</p>
</body>
</html>',
    'Organization Invitation

Hi {{.Name}},

You have been invited to join {{.OrganizationName}} on Booking System.

Role: {{.Role}}

Click the link below to accept the invitation:
{{.InvitationURL}}

This invitation will expire in {{.ExpiryDays}} days.

Best regards,
Booking System Team',
    '{"Name": "string", "OrganizationName": "string", "Role": "string", "InvitationURL": "string", "ExpiryDays": "number"}'
);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers to automatically update updated_at
CREATE TRIGGER update_email_jobs_updated_at 
    BEFORE UPDATE ON email_jobs 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_email_templates_updated_at 
    BEFORE UPDATE ON email_templates 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column(); 