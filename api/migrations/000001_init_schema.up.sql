-- =========================================================
-- USERS
-- =========================================================
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    email           VARCHAR(255) NOT NULL UNIQUE,
    username        VARCHAR(100) NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    auth_provider   VARCHAR(50) DEFAULT 'local',
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email
ON users(email);

CREATE INDEX idx_users_username
ON users(username);


-- =========================================================
-- ROLES
-- =========================================================

CREATE TABLE roles (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(50) NOT NULL UNIQUE,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_roles_name
ON roles(name);



-- =========================================================
-- USER ROLES
-- =========================================================

CREATE TABLE user_roles (
    user_id         BIGINT NOT NULL,
    role_id         BIGINT NOT NULL,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, role_id),

    CONSTRAINT fk_user_roles_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_roles_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_user_roles_user_id
ON user_roles(user_id);

CREATE INDEX idx_user_roles_role_id
ON user_roles(role_id);



-- =========================================================
-- ENVIRONMENTS
-- =========================================================

CREATE TABLE environments (
    id                  BIGSERIAL PRIMARY KEY,

    name                VARCHAR(255) NOT NULL,
    type                VARCHAR(50) NOT NULL,
    endpoint_url        TEXT NOT NULL,
    description         TEXT,
    status              VARCHAR(50) NOT NULL DEFAULT 'offline',
    tls_enabled         BOOLEAN NOT NULL DEFAULT FALSE,
    created_by          BIGINT,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_environments_created_by
        FOREIGN KEY (created_by)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_environments_name
ON environments(name);

CREATE INDEX idx_environments_type
ON environments(type);

CREATE INDEX idx_environments_status
ON environments(status);



-- =========================================================
-- STACKS
-- =========================================================

CREATE TABLE stacks (
    id                  BIGSERIAL PRIMARY KEY,
    environment_id      BIGINT NOT NULL,
    name                VARCHAR(255) NOT NULL,
    compose_content     TEXT NOT NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'created',
    created_by          BIGINT,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_stacks_environment
        FOREIGN KEY (environment_id)
        REFERENCES environments(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_stacks_created_by
        FOREIGN KEY (created_by)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_stacks_environment_id
ON stacks(environment_id);

CREATE INDEX idx_stacks_name
ON stacks(name);

CREATE INDEX idx_stacks_status
ON stacks(status);



-- =========================================================
-- AUDIT LOGS
-- =========================================================

CREATE TABLE audit_logs (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT,
    action              VARCHAR(255) NOT NULL,
    resource_type       VARCHAR(100) NOT NULL,
    resource_id         VARCHAR(255),
    ip_address          INET,
    metadata_json       JSONB,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_audit_logs_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_audit_logs_user_id
ON audit_logs(user_id);

CREATE INDEX idx_audit_logs_action
ON audit_logs(action);

CREATE INDEX idx_audit_logs_resource_type
ON audit_logs(resource_type);

CREATE INDEX idx_audit_logs_created_at
ON audit_logs(created_at);



-- =========================================================
-- DEFAULT ROLES
-- =========================================================

INSERT INTO roles (name, description)
VALUES
    ('admin', 'Full system access'),
    ('operator', 'Can manage environments and stacks'),
    ('viewer', 'Read-only access');





