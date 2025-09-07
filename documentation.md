# AffPilot Auth Service - Documentation

## Table of Contents
1. [Overview](#overview)
2. [Setup and Installation](#setup-and-installation)
   - [Prerequisites](#prerequisites)
   - [Local Setup](#local-setup)
   - [Docker Setup](#docker-setup)
3. [Environment Variables](#environment-variables)
4. [Features](#features)
5. [Technology Stack](#technology-stack)
6. [Project Structure](#project-structure)
5. [Database Schema](#database-schema)
6. [API Documentation](#api-documentation)

9. [User Types and Permissions](#user-types-and-permissions)
10. [Authentication Flow](#authentication-flow)
11. [Email Verification Process](#email-verification-process)
12. [Account Deletion Workflow](#account-deletion-workflow)
13. [Troubleshooting](#troubleshooting)

## Overview

AffPilot Auth Service is a robust authentication and authorization service for managing user identities, issuing JWT tokens, and implementing role-based access controls. The service provides secure user authentication, comprehensive identity management, and role-based authorization for applications.

### Local Setup

1. Clone the repository
```bash
git clone https://github.com/soyaibzihad10/Developer-Assignment
cd Developer-Assignment
go mod tidy
```

2. Configure environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

2. Start the application with Docker Compose
```bash
docker-compose up --build
```

This command will:
- Build the application Docker image
- Start the PostgreSQL database
- Apply the migrations
- Start the application on the configured port (default: 8080)

## Environment Variables

| Variable                | Description                                  | Example Value                       |
|-------------------------|----------------------------------------------|-------------------------------------|
| APP_ENV                 | Application environment                      | development                         |
| SERVER_PORT             | HTTP server port                             | 8080                                |
| LOG_LEVEL               | Logging level                                | debug                               |
| DB_HOST                 | PostgreSQL host                              | localhost                           |
| DB_PORT                 | PostgreSQL port                              | 5432                                |
| DB_USER                 | PostgreSQL username                          | postgres                            |
| DB_PASSWORD             | PostgreSQL password                          | postgres                            |
| DB_NAME                 | PostgreSQL database name                     | affpilot_auth                       |
| JWT_SECRET              | Secret key for signing JWT tokens            | your-secret-key-here                |
| JWT_EXPIRY              | JWT token validity duration                  | 24h                                 |
| SYSTEM_ADMIN_USERNAME   | Initial system admin username                | admin                               |
| SYSTEM_ADMIN_PASSWORD   | Initial system admin password                | adminpassword                       |
| SYSTEM_ADMIN_EMAIL      | Initial system admin email                   | admin@example.com                   |
| PASSWORD_SALT           | Salt for password hashing                    | your-password-salt-here             |
| EMAIL_VERIFICATION_URL  | Base URL for email verification links        | http://localhost:8080/api/v1/auth/verify |
| EMAIL_FROM              | Sender email address for system emails       | no-reply@example.com                |
| EMAIL_HOST              | SMTP server host                             | smtp.example.com                    |
| EMAIL_PORT              | SMTP server port                             | 587                                 |
| EMAIL_USERNAME          | SMTP server username                         | smtp-user                           |
| EMAIL_PASSWORD          | SMTP server password                         | smtp-password                       |
| EMAIL_SECURE            | Use TLS for SMTP (true/false)                | true                                |
| VERIFICATION_TOKEN_TTL  | Verification token lifetime in minutes       | 5                                   |


## Features

### Authentication
- Username/password login
- Secure password hashing with bcrypt
- Session management with JWT
- Email verification with time-limited hash links
- Password reset functionality

### Authorization
- Role-based permissions system with four user types:
  - System Admin: Full system access (defined at initialization)
  - Admin: Full access except cannot modify other admins
  - Moderator: Can manage all data but cannot delete users
  - User: Can only manage their own data
- Permission-based access control

### User Management
- Create, read, update user accounts
- Role assignment (admin only)
- Email verification workflow
- User deletion request workflow

### Administration
- Initial system admin credentials defined in environment variables
- Protected admin functions
- Audit trail for role assignments
- Tiered administrative privileges

## Technology Stack

| Component    | Technology    |
|--------------|---------------|
| Language     | Go (Golang)   |
| Database     | PostgreSQL    |
| Security     | JWT, bcrypt   |
| Deployment   | Docker        |

## Project Structure

```
.
├── cmd/
│   └── main.go
├── deployments/
├── internal/
│   ├── config/         # App configurations
│   ├── database/       # Database connections, migrations
│   ├── models/         # Data models
│   ├── services/       # Business logic
│   └── http/
│       ├── handlers/   # HTTP handlers
│       ├── middleware/ # Middleware components
│       └── routes/     # API route definitions
├── migrations/         # SQL migration files
├── scripts/            # Utility scripts
├── static/             # Static assets
├── test/               # Test cases and mocks
├── .env.example        # Environment variables template
├── docker-compose.yml  # Docker compose configuration
├── Dockerfile          # Docker build instructions
├── README.md
├── go.mod
└── go.sum
```

## Database Schema

The service uses PostgreSQL as its database with the following schema:

### Users Table

| Column              | Type        | Constraints              | Description                        |
|---------------------|-------------|---------------------------|------------------------------------|
| id                  | uuid        | PRIMARY KEY              | Unique user identifier             |
| username            | varchar(50) | UNIQUE, NOT NULL         | User's login name                  |
| email               | varchar(100)| UNIQUE, NOT NULL         | User's email address               |
| password_hash       | varchar(100)| NOT NULL                 | bcrypt hashed password             |
| first_name          | varchar(50) | NULL                     | User's first name                  |
| last_name           | varchar(50) | NULL                     | User's last name                   |
| email_verified      | boolean     | DEFAULT false            | Email verification status          |
| user_type           | varchar(20) | NOT NULL                 | User type (User/Moderator/Admin/SystemAdmin) |
| verification_token  | varchar(100)| NULL                     | Email verification token           |
| token_expiry        | timestamp   | NULL                     | Verification token expiry time     |
| reset_token         | varchar(255)| NULL                     | Password reset token               |
| reset_token_expiry  | timestamp   | NULL                     | Password reset token expiry time   |
| deletion_requested  | boolean     | DEFAULT false            | User has requested account deletion|
| active              | boolean     | DEFAULT true             | Account active status              |
| created_at          | timestamp   | NOT NULL                 | Creation timestamp                 |
| updated_at          | timestamp   | NOT NULL                 | Last update timestamp              |

### Roles Table

| Column      | Type        | Constraints      | Description           |
|-------------|-------------|------------------|-----------------------|
| id          | uuid        | PRIMARY KEY      | Role identifier       |
| name        | varchar(50) | UNIQUE, NOT NULL | Role name             |
| description | text        | NULL             | Role description      |
| created_at  | timestamp   | NOT NULL         | Creation timestamp    |
| updated_at  | timestamp   | NOT NULL         | Last update timestamp |

### Permissions Table

| Column      | Type         | Constraints      | Description                |
|-------------|--------------|------------------|----------------------------|
| id          | uuid         | PRIMARY KEY      | Permission identifier      |
| name        | varchar(100) | UNIQUE, NOT NULL | Permission name            |
| description | text         | NULL             | Permission description     |
| resource    | varchar(100) | NOT NULL         | Resource being accessed    |
| action      | varchar(50)  | NOT NULL         | Action on resource (read, write) |
| created_at  | timestamp    | NOT NULL         | Creation timestamp         |
| updated_at  | timestamp    | NOT NULL         | Last update timestamp      |

### UserRoles Table

| Column      | Type      | Constraints                    | Description                |
|-------------|-----------|--------------------------------|----------------------------|
| user_id     | uuid      | FK -> users.id, NOT NULL       | Reference to user          |
| role_id     | uuid      | FK -> roles.id, NOT NULL       | Reference to role          |
| assigned_by | uuid      | FK -> users.id, NOT NULL       | Admin who assigned the role|
| created_at  | timestamp | NOT NULL                       | Assignment timestamp       |
| PRIMARY KEY | (user_id, role_id)                         | Composite primary key      |

### RolePermissions Table

| Column        | Type      | Constraints                        | Description           |
|---------------|-----------|------------------------------------|-----------------------|
| role_id       | uuid      | FK -> roles.id, NOT NULL           | Reference to role     |
| permission_id | uuid      | FK -> permissions.id, NOT NULL     | Reference to permission|
| created_at    | timestamp | NOT NULL                           | Assignment timestamp  |
| PRIMARY KEY   | (role_id, permission_id)                       | Composite primary key |

## API Documentation

### Authentication Endpoints

| Method | Endpoint                        | Description                  | Access Level  |
|--------|----------------------------------|------------------------------|---------------|
| POST   | /api/v1/auth/login               | Authenticate user            | Public        |
| POST   | /api/v1/auth/logout              | Invalidate session           | Authenticated |
| POST   | /api/v1/auth/register            | Register a user              | Public        |
| GET    | /api/v1/auth/verify              | Verify email with token      | Public        |
| POST   | /api/v1/auth/resend-verification | Resend verification email    | Public        |
| POST   | /api/v1/auth/password-reset      | Reset password               | Public        |

### User Management Endpoints

| Method | Endpoint                             | Description                | Access Level       |
|--------|--------------------------------------|----------------------------|-------------------|
| GET    | /api/v1/users                        | List all users             | Admin+            |
| GET    | /api/v1/users/{user_id}              | Get user details           | User (self only)/Moderator+ |
| PUT    | /api/v1/users/{user_id}              | Update user                | User (self only)/Admin+ |
| POST   | /api/v1/users/{user_id}/request-deletion | Request account deletion  | User (self only)  |
| DELETE | /api/v1/users/{user_id}              | Delete user                | Moderator+        |

### Role Management Endpoints

| Method | Endpoint                             | Description                | Access Level       |
|--------|--------------------------------------|----------------------------|-------------------|
| POST   | /api/v1/users/{user_id}/role         | Change user role           | Admin+            |
| POST   | /api/v1/users/{user_id}/promote/admin | Promote to admin          | System Admin      |
| POST   | /api/v1/users/{user_id}/promote/moderator | Promote to moderator  | Admin+            |
| POST   | /api/v1/users/{user_id}/demote       | Demote user role           | Admin+            |
| GET    | /api/v1/roles                        | List all roles             | Admin+            |
| GET    | /api/v1/roles/{role_id}              | Get role details           | Admin+            |
| POST   | /api/v1/roles                        | Create role                | Admin+            |
| PUT    | /api/v1/roles/{role_id}              | Update role                | Admin+            |
| DELETE | /api/v1/roles/{role_id}              | Delete role                | Admin+            |

### Permission Endpoints

| Method | Endpoint                             | Description                | Access Level       |
|--------|--------------------------------------|----------------------------|-------------------|
| GET    | /api/v1/permissions                  | List all permissions       | Admin+            |
| GET    | /api/v1/permissions/{permission_id}  | Get permission details     | Admin+            |

### User Profile Endpoints

| Method | Endpoint                             | Description                | Access Level       |
|--------|--------------------------------------|----------------------------|-------------------|
| GET    | /api/v1/me                           | Get current user profile   | Authenticated     |
| GET    | /api/v1/me/permissions               | Get current user permissions | Authenticated   |

## Setup and Installation

### Prerequisites

- Go 1.20+
- Docker and Docker Compose
- PostgreSQL 13+
- Git


## User Types and Permissions

### System Admin
- Created when the database is initialized (from environment variables)
- Has complete access to all system functions
- Cannot be deleted from the system
- Can promote users to any role
- Can demote any user except other System Admins
- Can manage all roles and permissions

### Admin
- Has full administrative access
- Cannot remove System Admins or other Admins
- Cannot promote users to System Admin role
- Can promote users to Moderator or Admin
- Can manage all roles and permissions

### Moderator
- Can create, read, update, and delete data
- Cannot delete user accounts directly
- Can process deletion requests from users
- Can manage content but not system configuration

### User
- Can only create, read, update, and delete their own data
- Cannot access other users' data
- Can request account deletion (to be processed by Moderator+)

## Authentication Flow

1. **Registration**:
   - User submits registration data (username, email, password)
   - System validates input
   - Password is hashed using bcrypt
   - User account is created in `users` table with `email_verified = false`
   - Verification email is sent

2. **Login**:
   - User submits username/email and password
   - System validates credentials
   - If valid, JWT token is issued containing user ID and role
   - Token is returned to client

3. **Session Management**:
   - Client includes JWT token in Authorization header
   - Server validates token on each request
   - Tokens expire based on JWT_EXPIRY configuration

## Email Verification Process

1. When a user registers:
   - Account is created with `email_verified = false`
   - A verification token is generated and stored in the user record
   - Token expiry is set to 5 minutes from creation time
   - Verification email is sent with link containing the token

2. When user clicks the verification link:
   - Token is validated against database
   - If valid and not expired, user's email is marked as verified
   - Verification token is cleared from database

3. If token has expired:
   - User is prompted to request a new verification email
   - New token is generated with new 5-minute expiry

## Account Deletion Workflow

1. User requests account deletion via API
   - User record is flagged with `deletion_requested = true`

2. Moderator, Admin, or System Admin reviews the request
   - Reviews can be done through the admin interface

3. Upon approval:
   - User account is permanently deleted from the system
   - All associated data is removed


### API Testing with Postman

A Postman collection is available for testing the API endpoints. Import the collection from:
[AffPilot Auth Service Postman Collection](https://none11-5627.postman.co/workspace/none-Workspace~97085782-b278-41b0-8928-215f543210b2/collection/37975347-721220c1-0c48-47b9-a9d4-841299577403?action=share&creator=37975347)

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   - Ensure PostgreSQL is running
   - Verify database credentials in .env file
   - Check that the database exists and is accessible

2. **Email Sending Failures**
   - Verify SMTP credentials in .env file
   - Check that the SMTP server is accessible from your network
   - Look for email sending errors in application logs

3. **JWT Token Issues**
   - Ensure JWT_SECRET is properly set
   - Check token expiration settings
   - Verify that client is sending token correctly in Authorization header

4. **Migration Failures**
   - Ensure database is accessible
   - Check that migration files are in the correct order
   - Look for SQL syntax errors in migration files

### Logging

The application logs can be helpful for diagnosing issues:

```bash
# View application logs when running with Docker
docker-compose logs app

# Increase log verbosity
# Set LOG_LEVEL=debug in .env file
```

For any other issues or questions, please open an issue on the repository or contact the development team.
