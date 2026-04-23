DROP TABLE IF EXISTS users;

CREATE TABLE users
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    login      VARCHAR     NOT NULL UNIQUE,
    name       VARCHAR     NOT NULL,
    password   VARCHAR     NOT NULL,
    admin      BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO users (login, name, password, admin)
VALUES ('admin', 'admin', '$2a$10$RWJ5rvX1ZCxOGs1E7npb6Od.cjZ/eF2FvdYGxBqJhb23DHSF4dJyG', true);


CREATE INDEX idx_users_admin ON users (admin);

CREATE INDEX idx_users_created_at ON users (created_at);