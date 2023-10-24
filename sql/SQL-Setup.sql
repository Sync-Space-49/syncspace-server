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
    description     TEXT
);

CREATE TABLE IF NOT EXISTS Boards (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL, -- Changed to Title from Name to match Lists and Cards
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- 'default' calling a function may not work? if not j remove the function call
    modified_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_private      BOOLEAN DEFAULT FALSE,                  -- defaults to public
    organization_id UUID, FOREIGN KEY (organization_id) REFERENCES Organizations(id)
);

CREATE TABLE IF NOT EXISTS Panels ( -- Changed to Panels from Lists
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    position        SMALLINT, -- changed from int to smallint (will not have more than 32767 lists/cards/etc)
    board_id        SERIAL, FOREIGN KEY (board_id) REFERENCES Boards(id)
);

CREATE TABLE IF NOT EXISTS Stacks (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    position        SMALLINT, -- changed from int to smallint (will not have more than 32767 lists/cards/etc)
    panel_id        SERIAL, FOREIGN KEY (panel_id) REFERENCES Panels(id)
);

CREATE TABLE IF NOT EXISTS Cards (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    position        SMALLINT, -- changed from int to smallint (will not have more than 32767 lists/cards/etc)
    stack_id         SERIAL, FOREIGN KEY (stack_id) REFERENCES Stacks(id)
);

CREATE TABLE IF NOT EXISTS Assigned_Cards (
    id              SERIAL PRIMARY KEY, -- changed from mere join table structure to give an ID for future ease of reference
    user_id         VARCHAR(64),
    card_id         SERIAL, FOREIGN KEY (card_id) REFERENCES Cards(id)
);

CREATE TABLE IF NOT EXISTS Tags (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Card_Tags (
    id              SERIAL PRIMARY KEY, -- changed from mere join table structure to give an ID for future ease of reference
    tag_id          SERIAL, FOREIGN KEY (tag_id) REFERENCES Tags(id),
    card_id         SERIAL, FOREIGN KEY (card_id) REFERENCES Cards(id)
);