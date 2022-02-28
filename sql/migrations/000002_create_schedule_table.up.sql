CREATE TABLE IF NOT EXISTS "schedule"(
    "id" VARCHAR(50) PRIMARY KEY,
    "event_id" VARCHAR(50) NOT NULL,
    "start_time" BIGINT NOT NULL,
    "duration" BIGINT NOT NULL,
    "is_full_day" BOOLEAN NOT NULL,
    "recurring_interval" BIGINT NOT NULL DEFAULT 0,
    CONSTRAINT "fk_event" FOREIGN KEY ("event_id") REFERENCES event("id") ON DELETE CASCADE
);