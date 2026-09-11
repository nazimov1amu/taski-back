CREATE TABLE tasks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    project_id  UUID REFERENCES projects (id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    description TEXT,
    completed   BOOLEAN NOT NULL DEFAULT false,
    planned_at  TEXT,
    start_time  TEXT,
    end_time    TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_user_id ON tasks (user_id);
CREATE INDEX idx_tasks_project_id ON tasks (project_id);
CREATE INDEX idx_tasks_planned_at ON tasks (planned_at);
