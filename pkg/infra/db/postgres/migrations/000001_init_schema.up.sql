CREATE TYPE task_status AS ENUM ('pending', 'processing', 'failed', 'dlq');

CREATE TABLE IF NOT EXISTS tasks (
                                     id TEXT PRIMARY KEY,
                                     channel TEXT NOT NULL,
                                     payload BYTEA NOT NULL,
                                     metadata JSONB DEFAULT '{}',
                                     status task_status DEFAULT 'pending',
                                     attempts INT DEFAULT 0,
                                     created_at TIMESTAMPTZ DEFAULT NOW(),
    visible_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_tasks_polling ON tasks (channel, visible_at, status)
    WHERE status = 'pending';