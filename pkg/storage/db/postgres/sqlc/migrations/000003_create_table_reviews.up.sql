DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS tasks;

CREATE TYPE review_status AS ENUM ('pending', 'processing', 'failed', 'done');


CREATE TABLE IF NOT EXISTS tasks
(
    id   int PRIMARY KEY,
    name varchar
);


CREATE TABLE IF NOT EXISTS reviews
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    userid     BIGINT,
    task       int not null,
    status     review_status DEFAULT 'pending',
    attempts   INT           DEFAULT 0,
    created_at TIMESTAMPTZ   DEFAULT NOW(),
    fileID     varchar,
    CONSTRAINT fk_reviews_user
        FOREIGN KEY (userid)
            REFERENCES users (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_reviews_task
        FOREIGN KEY (task)
            REFERENCES tasks (id)

);

CREATE INDEX idx_reviews_user_id ON reviews (userid);
CREATE INDEX idx_reviews_task_id ON reviews (task);
CREATE INDEX idx_reviews_status ON reviews (status);