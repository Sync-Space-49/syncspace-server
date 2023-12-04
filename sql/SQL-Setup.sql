-- SyncSpace SQL Setup

-- Created for https://github.com/Sync-Space-49/syncspace-server/

-- Authors:
--     Nathan Laney
--     Dylan Halstead

SET TIMEZONE = "America/New_York"; -- sets timezone for timestamptz datatypes

CREATE TABLE IF NOT EXISTS Organizations (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    owner_id        VARCHAR(64) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT, 
    ai_enabled      BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS Boards (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     TEXT DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_private      BOOLEAN DEFAULT FALSE,                  -- defaults to public
    organization_id UUID, FOREIGN KEY (organization_id) REFERENCES Organizations(id) ON DELETE CASCADE,
    owner_id        VARCHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS Panels (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    position        SMALLINT,
    board_id        UUID, FOREIGN KEY (board_id) REFERENCES Boards(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Stacks (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    position        SMALLINT,
    panel_id        UUID, FOREIGN KEY (panel_id) REFERENCES Panels(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Cards (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    points          VARCHAR(30) NOT NULL DEFAULT 0;
    position        SMALLINT,
    stack_id        UUID, FOREIGN KEY (stack_id) REFERENCES Stacks(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Assigned_Cards (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id         VARCHAR(64), 
    card_id         UUID, FOREIGN KEY (card_id) REFERENCES Cards(id)
);

CREATE TABLE IF NOT EXISTS Tags (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name            VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Card_Tags (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    tag_id          UUID, FOREIGN KEY (tag_id) REFERENCES Tags(id),
    card_id         UUID, FOREIGN KEY (card_id) REFERENCES Cards(id)
);