CREATE TABLE IF NOT EXISTS "event"(
    "id" VARCHAR(50) PRIMARY KEY,
    "title" VARCHAR(200) NOT NULL DEFAULT '',
    "description" TEXT NOT NULL DEFAULT '',
    "timezone" VARCHAR(30) NOT NULL,
    "created_by" VARCHAR(50) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NULL
);