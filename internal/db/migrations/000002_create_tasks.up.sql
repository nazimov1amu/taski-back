CREATE TABLE tasks (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  uuid REFERENCES projects (id) ON DELETE SET NULL,
    title       text NOT NULL,
    description text NOT NULL DEFAULT '',
    completed   boolean NOT NULL DEFAULT false,
    start_time  timestamp without time zone NOT NULL DEFAULT now(),
    end_time    timestamp without time zone NOT NULL DEFAULT now(),
    created_at  timestamp NOT NULL DEFAULT now(),
    updated_at  timestamp NOT NULL DEFAULT now()
);

