CREATE TABLE IF NOT EXISTS message (
    "id"             uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "room_id"        uuid             NOT NULL,
    "text"           VARCHAR(255)     NOT NULL,
    "reaction_count" BIGINT           NOT NULL DEFAULT 0,
    "answered"       BOOLEAN          NOT NULL DEFAULT false,
    "created_at"     TIMESTAMP        NOT NULL DEFAULT NOW(),
    "updated_at"     TIMESTAMP        NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (room_id) REFERENCES room (id) ON DELETE CASCADE
);

---- create above / drop below ----

DROP TABLE IF EXISTS message;
