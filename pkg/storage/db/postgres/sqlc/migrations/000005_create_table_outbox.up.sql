DROP TABLE IF EXISTS outbox_reviews;

CREATE TYPE outbox_status AS ENUM ('pending', 'processing', 'failed', 'done');

CREATE TABLE outbox_reviews
(
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_id      UUID  NOT NULL UNIQUE,
    userid        BIGINT,
    reviewid      BIGINT,
    payload       JSONB NOT NULL,
    status        outbox_status  DEFAULT 'pending' NOT NULL,
    created_at    TIMESTAMPTZ    DEFAULT NOW(),
    processed_at  TIMESTAMPTZ,
    attempts      INT            DEFAULT 0,
    max_attempts  INT   NOT NULL DEFAULT 5,
    next_retry_at TIMESTAMPTZ,
    last_error    TEXT
);

CREATE INDEX idx_outbox_reviews_polling
    ON outbox_reviews (status, next_retry_at, created_at);

CREATE INDEX idx_outbox_reviews_user
    ON outbox_reviews (userid);

CREATE INDEX idx_outbox_reviews_review
    ON outbox_reviews (reviewid);

CREATE INDEX idx_outbox_reviews_poll
    ON outbox_reviews (status, next_retry_at, created_at)