CREATE TABLE IF NOT EXISTS users (
    user_id     VARCHAR(64) PRIMARY KEY,
    passwd_hash TEXT        NOT NULL,
    create_at   TIMESTAMP   NOT NULL DEFAULT(now()),
    update_at   TIMESTAMP   NOT NULL DEFAULT(now())
);

CREATE TABLE IF NOT EXISTS passwords (
    id          SERIAL       PRIMARY KEY,
    user_id     VARCHAR(64)  NOT NULL,
    name	    VARCHAR(128) NOT NULL,			
    username    BYTEA        NOT NULL,
    password    BYTEA        NOT NULL,
    notes       BYTEA        NOT NULL,
    create_at   TIMESTAMP    NOT NULL DEFAULT(now()),
    update_at   TIMESTAMP    NOT NULL DEFAULT(now())
);
CREATE INDEX IF NOT EXISTS passwords_user_id_idx 
ON passwords (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS passwords_user_id_name_idx 
ON passwords (user_id, name);	

CREATE TABLE IF NOT EXISTS notes (
    id        SERIAL       PRIMARY KEY,
    user_id   VARCHAR(64)  NOT NULL,
    name	  VARCHAR(128) NOT NULL,
    notes     BYTEA        NOT NULL,
    create_at TIMESTAMP    NOT NULL DEFAULT(now()),
    update_at TIMESTAMP    NOT NULL DEFAULT(now())
);
CREATE INDEX IF NOT EXISTS notes_user_id_idx 
ON notes (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS notes_user_id_name_idx 
ON notes (user_id, name);	

CREATE TABLE IF NOT EXISTS cards (
    id          SERIAL       PRIMARY KEY,
    user_id     VARCHAR(64)  NOT NULL,
    name        VARCHAR(128) NOT NULL,
    number      BYTEA		 NOT NULL,
    pin         BYTEA        NOT NULL,	 
    notes       BYTEA        NOT NULL,
    create_at   TIMESTAMP    NOT NULL DEFAULT(now()),
    update_at   TIMESTAMP    NOT NULL DEFAULT(now())
);
CREATE INDEX IF NOT EXISTS cards_user_id_idx 
ON cards (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS cards_user_id_name_idx 
ON cards (user_id, name);

CREATE TABLE IF NOT EXISTS binaries (
    id        SERIAL       PRIMARY KEY,
    user_id   VARCHAR(64)  NOT NULL,
    name      VARCHAR(128) NOT NULL,
    size      INT		   NOT NULL,
    notes     BYTEA        NOT NULL,
    bin_id    OID		   NOT NULL,
    create_at TIMESTAMP    NOT NULL DEFAULT(now()),
    update_at TIMESTAMP    NOT NULL DEFAULT(now())
);
CREATE INDEX IF NOT EXISTS binaries_user_id_idx 
ON binaries (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS binaries_user_id_name_idx
ON binaries (user_id, name);