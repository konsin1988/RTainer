CREATE TABLE teams (
    id                  BIGSERIAL PRIMARY KEY,
    name                VARCHAR(255) NOT NULL UNIQUE,
    keycloak_group_id   VARCHAR(255) NOT NULL UNIQUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE environment_permissions (
    id                  BIGSERIAL PRIMARY KEY,
    environment_id      BIGINT NOT NULL,
    team_id             BIGINT NOT NULL,
    permission          VARCHAR(50) NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_env_permissions_environment
        FOREIGN KEY (environment_id)
        REFERENCES environments(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_env_permissions_team
        FOREIGN KEY (team_id)
        REFERENCES teams(id)
        ON DELETE CASCADE
);
