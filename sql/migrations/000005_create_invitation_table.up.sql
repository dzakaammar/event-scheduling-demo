CREATE TABLE IF NOT EXISTS "event_invitation"(
    "id" VARCHAR(50) PRIMARY KEY,
    "event_id" VARCHAR(50) NOT NULL,
    "user_id" INTEGER NOT NULL,
    "token" TEXT NOT NULL UNIQUE,
    "status" SMALLINT NOT NULL,
    "updated_at" TIMESTAMP NULL,
    CONSTRAINT "fk_event" FOREIGN KEY ("event_id") REFERENCES event("id") ON DELETE CASCADE
);