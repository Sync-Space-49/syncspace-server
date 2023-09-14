-- SyncSpace SQL Setup

-- Created for https://github.com/Sync-Space-49/syncspace-server/

-- Authors: 
--     Nathan Laney
--     Dylan Halstead

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
    description     TEXT
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
    is_private      BOOLEAN DEFAULT FALSE,                  -- defaults to public
    organization_id SERIAL, FOREIGN KEY (organization_id) REFERENCES Organizations(id)
);

CREATE TABLE IF NOT EXISTS Lists (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    position        SMALLINT, -- changed from int to smallint (will not have more than 32767 lists/cards/etc)
    board_id        SERIAL, FOREIGN KEY (board_id) REFERENCES Boards(id)
);

CREATE TABLE IF NOT EXISTS Cards (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    position        SMALLINT, -- changed from int to smallint (will not have more than 32767 lists/cards/etc)
    list_id         SERIAL, FOREIGN KEY (list_id) REFERENCES Lists(id)
);

CREATE TABLE IF NOT EXISTS Assigned_Cards (
    id              SERIAL PRIMARY KEY, -- changed from mere join table structure to give an ID for future ease of reference
    user_id         SERIAL, FOREIGN KEY (user_id) REFERENCES Users(id),
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

CREATE TABLE IF NOT EXISTS Board_Members (
    id              SERIAL PRIMARY KEY,
    member_id       SERIAL, FOREIGN KEY (member_id) REFERENCES Organization_Members(id),
    board_id        SERIAL, FOREIGN KEY (board_id) REFERENCES Boards(id)
);

CREATE TABLE IF NOT EXISTS Board_Roles (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    is_default      BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS Board_Member_Roles (
    id              SERIAL PRIMARY KEY,
    board_member_id SERIAL, FOREIGN KEY (board_member_id) REFERENCES Board_Members(id),
    board_role_id   SERIAL, FOREIGN KEY (board_role_id) REFERENCES Board_Roles(id)
);

CREATE TABLE IF NOT EXISTS Board_Privileges (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    description     TEXT
);

CREATE TABLE IF NOT EXISTS Board_Roles_Privileges (
    id              SERIAL PRIMARY KEY,
    board_privilege_id
                    SERIAL, FOREIGN KEY (board_privilege_id) REFERENCES Board_Privileges(id),
    board_role_id   SERIAL, FOREIGN KEY (board_role_id) REFERENCES Board_Roles(id)
);

CREATE TABLE IF NOT EXISTS Organization_Roles (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    is_default      BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS Organization_Privileges (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    description     TEXT
);

CREATE TABLE IF NOT EXISTS Organization_Member_Roles (
    id              SERIAL PRIMARY KEY,
    organization_role_id
                    SERIAL, FOREIGN KEY (organization_role_id) REFERENCES Organization_Roles(id),
    organization_member_id
                    SERIAL, FOREIGN KEY (organization_member_id) REFERENCES Organization_Members(id)
);


CREATE TABLE IF NOT EXISTS Organization_Role_Priviliges (
    id              SERIAL PRIMARY KEY,
    organization_role_id
                    SERIAL, FOREIGN KEY (organization_role_id) REFERENCES Organization_Roles(id),
    organization_privilege_id
                    SERIAL, FOREIGN KEY (organization_privilege_id) REFERENCES Organization_Privileges(id)
);
