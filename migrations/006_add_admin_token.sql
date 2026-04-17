-- Add admin_token column to provider_applications
ALTER TABLE provider_applications ADD COLUMN admin_token VARCHAR(64) UNIQUE;

-- Create index for faster token lookups
CREATE INDEX idx_provider_applications_admin_token ON provider_applications(admin_token);