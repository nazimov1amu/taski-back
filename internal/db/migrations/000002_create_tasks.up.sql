CREATE TABLE tasks (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  uuid REFERENCES projects (id) ON DELETE CASCADE,
    title       text NOT NULL,
    description text DEFAULT '',
    completed   boolean NOT NULL DEFAULT false,
    planned_at  date,
    start_time  timestamp without time zone DEFAULT now(),
    end_time    timestamp without time zone DEFAULT now(),
    created_at  timestamp NOT NULL DEFAULT now(),
    updated_at  timestamp NOT NULL DEFAULT now()
);
