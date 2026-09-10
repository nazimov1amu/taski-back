CREATE TYPE sync_event_type AS ENUM ('create', 'update', 'delete');

CREATE TABLE sync_events(
    id          BIGSERIAL PRIMARY KEY NOT NULL,
    user_id     UUID NOT NULL REFERENCES users(id),
    sequence_id BIGINT NOT NULL,
    entity_type VARCHAR(255) NOT NULL,
    entity_id   UUID NOT NULL,
    operation   sync_event_type NOT NULL,
    payload     JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
