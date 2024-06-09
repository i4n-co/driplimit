-- create keyspaces table
CREATE TABLE
    IF NOT EXISTS keyspaces (
        ksid text PRIMARY KEY,
        name text NOT NULL,
        keys_prefix text NOT NULL,
        rate_limit_limit int DEFAULT 0,
        rate_limit_refill_rate int DEFAULT 0,
        rate_limit_refill_interval int DEFAULT 0,
        deleted_at INT NOT NULL DEFAULT 0
    );

CREATE INDEX idx_keyspaces_deleted_at ON keyspaces (deleted_at);
CREATE UNIQUE INDEX idx_unique_keyspaces_name ON keyspaces (name) WHERE deleted_at = 0;
CREATE UNIQUE INDEX idx_unique_keyspaces_keys_prefix ON keyspaces (keys_prefix) WHERE deleted_at = 0;
-- Create a view to filter out deleted keyspaces
CREATE VIEW v_keyspaces AS SELECT * FROM keyspaces WHERE deleted_at = 0;


-- create table keys
CREATE TABLE
    IF NOT EXISTS keys (
        kid text PRIMARY KEY,
        ksid text NOT NULL,
        last_used int DEFAULT 0,
        rate_limit_state_last_refilled int DEFAULT (unixepoch ('now')),
        rate_limit_state_remaining int DEFAULT 0,
        rate_limit_refill_rate int DEFAULT 0,
        rate_limit_refill_interval int DEFAULT 0,
        rate_limit_limit int DEFAULT 0,
        expires_at int DEFAULT (unixepoch ('now', '+99 year')),
        created_at int DEFAULT (unixepoch ('now')),
        token_hash text NOT NULL,
        deleted_at INT NOT NULL DEFAULT 0,
        FOREIGN KEY (ksid) REFERENCES keyspaces (ksid)
    );

-- create keys indices
CREATE INDEX IF NOT EXISTS idx_keys_token_hash ON keys (token_hash);
CREATE INDEX IF NOT EXISTS idx_keys_ksid ON keys (ksid);
-- Create a view to filter out deleted keys
CREATE VIEW v_keys AS SELECT * FROM keys WHERE deleted_at = 0;


-- create service keys table
CREATE TABLE service_keys (
    skid                TEXT NOT NULL,
    token_hash          TEXT NOT NULL,
    admin               INTEGER NOT NULL default 0,
    description         TEXT NOT NULL default '',
    created_at          INTEGER NOT NULL default  (unixepoch ('now')),
    deleted_at          INTEGER NOT NULL default  0,
    PRIMARY KEY (skid)
);

-- Create an index on the token_hash column
CREATE INDEX idx_service_keys_token_hash ON service_keys (token_hash);

-- Create an index on the deleted_at column
CREATE INDEX idx_service_keys_deleted_at ON service_keys (deleted_at);

-- Create a view to filter out deleted service keys
CREATE VIEW v_service_keys AS
    SELECT * FROM service_keys WHERE deleted_at = 0;

-- Create a table to store the keyspaces that a service key has access to
-- foreign key to ksid is not added because ksid can be '*' which cannot be a valid ksid
-- delete_at is not added as we prefer to hard delete policies
CREATE TABLE keyspaces_policies (
    skid                TEXT        NOT NULL,
    ksid                TEXT        NOT NULL,
    read                INTEGER     NOT NULL default 0,
    write               INTEGER     NOT NULL default 0,
    PRIMARY KEY (skid, ksid),
    FOREIGN KEY (skid) REFERENCES service_keys (skid)
);

