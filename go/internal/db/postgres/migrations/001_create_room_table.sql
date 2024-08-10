CREATE TABLE IF NOT EXISTS room (
    "id"          uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "name"        VARCHAR(255)     NOT NULL,
    "created_at"  TIMESTAMP        NOT NULL DEFAULT NOW(),
    "updated_at"  TIMESTAMP        NOT NULL DEFAULT NOW()
);

---- create above / drop below ----

DROP TABLE IF EXISTS room;
