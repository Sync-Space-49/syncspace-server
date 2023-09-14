```
SyncSpace SQL Setup

Created for https://github.com/Sync-Space-49/syncspace-server/

Authors: 
    Nathan Laney
    Dylan Halstead

https://user-images.githubusercontent.com/70990184/267409122-ddc40017-4d0a-4e29-a707-4858b364b630.png
``` 
SET TIMEZONE = "America/New_York"; -- sets timezone for timestamptz datatypes

CREATE TABLE IF NOT EXISTS Users (
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    password        VARCHAR(255) NOT NULL,
    pfp_url         VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS Organizations (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    -- desc length 255 too short? what limit
    description     VARCHAR(255) 
);

CREATE TABLE IF NOT EXISTS Organization_Members (
    id              SERIAL PRIMARY KEY,
    user_id         SERIAL, FOREIGN KEY (user_id) REFERENCES Users(id),
    organization_id SERIAL, FOREIGN KEY (organization_id) REFERENCES Organizations(id)
);

CREATE TABLE IF NOT EXISTS Boards (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL, -- Changed to Title from Name to match Lists and Cards 
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- 'default' calling a function may not work? if not j remove the function call
    modified_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_private      BOOLEAN DEFAULT 0;                  -- defaults to public
    organization_id SERIAL, FOREIGN KEY (organization_id) REFERENCES Organizations(id)
);

CREATE TABLE IF NOT EXISTS Lists (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    position        SMALLINT(255), -- changed from int to smallint (will not have more than 32767 lists/cards/etc)
    board_id        SERIAL, FOREIGN KEY (board_id) REFERENCES Boards(id)
);


CREATE TABLE IF NOT EXISTS Cards (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     VARCHAR(255),
    position        SMALLINT(255), -- changed from int to smallint (will not have more than 32767 lists/cards/etc)
    list_id         SERIAL, FOREIGN KEY (list_id) REFERENCES Lists(id)
);

CREATE TABLE IF NOT EXISTS Assigned_Cards (
    id              SERIAL PRIMARY KEY,
    user_id         SERIAL, FOREIGN KEY (user_id) REFERENCES Users(id),
    card_id         SERIAL, FOREIGN KEY (card_id) REFERENCES Cards(id)
);