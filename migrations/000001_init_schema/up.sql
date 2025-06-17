-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email_verified BOOLEAN DEFAULT FALSE,
    user_type VARCHAR(20) NOT NULL,
    verification_token VARCHAR(100),
    reset_token VARCHAR(100),
    token_expiry TIMESTAMP,
    deletion_requested BOOLEAN DEFAULT FALSE,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- UserRoles junction table
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

-- RolePermissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

-- Insert default roles
INSERT INTO roles (name, description) VALUES
    ('system_admin', 'Full system access with ability to manage all aspects of the system'),
    ('admin', 'Administrative access with limitations on managing other admins'),
    ('moderator', 'Can manage content but cannot delete users directly'),
    ('user', 'Regular user with access only to their own data');

-- Insert default permissions
INSERT INTO permissions (name, resource, action, description) VALUES
    ('user:read:all', 'user', 'read:all', 'Read all user data'),
    ('user:create:all', 'user', 'create:all', 'Create any user'),
    ('user:update:all', 'user', 'update:all', 'Update any user'),
    ('user:delete:all', 'user', 'delete:all', 'Delete any user'),
    ('user:read:self', 'user', 'read:self', 'Read own user data'),
    ('user:update:self', 'user', 'update:self', 'Update own user data'),
    ('user:delete:self', 'user', 'delete:self', 'Request own account deletion'),
    ('role:read', 'role', 'read', 'Read roles'),
    ('role:create', 'role', 'create', 'Create roles'),
    ('role:update', 'role', 'update', 'Update roles'),
    ('role:delete', 'role', 'delete', 'Delete roles'),
    ('permission:read', 'permission', 'read', 'Read permissions'),
    ('user:promote:admin', 'user', 'promote:admin', 'Promote user to admin'),
    ('user:promote:moderator', 'user', 'promote:moderator', 'Promote user to moderator'),
    ('user:demote', 'user', 'demote', 'Demote user role');

-- Assign permissions to roles
-- System Admin permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'system_admin'), 
    id 
FROM permissions;

-- Admin permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'admin'), 
    id 
FROM permissions
WHERE name NOT IN ('user:promote:admin');

-- Moderator permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'moderator'), 
    id 
FROM permissions
WHERE name IN (
    'user:read:all', 
    'user:read:self', 
    'user:update:self', 
    'user:delete:self',
    'user:delete:all'
);

-- User permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'user'), 
    id 
FROM permissions
WHERE name IN (
    'user:read:self', 
    'user:update:self', 
    'user:delete:self'
);
