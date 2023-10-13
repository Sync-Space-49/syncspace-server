-- SyncSpace SQL Filler Data Fill

-- Created for https://github.com/Sync-Space-49/syncspace-server/

-- Authors: 
--     Nathan Laney
--     Dylan Halstead

SET TIMEZONE = "America/New_York"; -- sets timezone for timestamptz datatypes

-- CREATE TABLE IF NOT EXISTS Organizations (
--     id              SERIAL PRIMARY KEY,
--     name            VARCHAR(255) NOT NULL,
--     description     TEXT
-- );

INSERT INTO Organizations(id, name, description) VALUES 
    (1, 'SyncSpace', 'This is the SyncSpace example organization description. This is a temporary description in filler data. '),
    (2, 'BioQuest', 'Lorem ipsum, dolor sit amet consectetur adipisicing elit. Molestias aut, repellat ipsum facere voluptate dicta obcaecati deserunt nobis suscipit eaque?'),
    (3, 'Some Recipe App', 'Lorem ipsum, vulputate eu scelerisque felis imperdiet proin fermentum leo vel orci porta non pulvinar neque laoreet.');

-- CREATE TABLE IF NOT EXISTS Boards (
--     id              SERIAL PRIMARY KEY,
--     title           VARCHAR(255) NOT NULL, -- Changed to Title from Name to match Lists and Cards
--     created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- 'default' calling a function may not work? if not j remove the function call
--     modified_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     is_private      BOOLEAN DEFAULT FALSE,                  -- defaults to public
--     organization_id SERIAL, FOREIGN KEY (organization_id) REFERENCES Organizations(id)
-- );

INSERT INTO Boards(id, title, organization_id) VALUES
    (1, 'SyncSpace Main Board', 1),
    (2, 'A second, cooler board', 1),
    (3, 'BioQuest Main Board', 2),
    (4, 'A second, cooler board', 2),
    (5, 'Board Title', 3),
    (6, 'An additional, even COOLER board', 1);

-- CREATE TABLE IF NOT EXISTS Panels ( -- Changed to Panels from Lists
--     id              SERIAL PRIMARY KEY,
--     title           VARCHAR(255) NOT NULL,
--     position        SMALLINT, -- changed from int to smallint (will not have more than 32767 lists/cards/etc)
--     board_id        SERIAL, FOREIGN KEY (board_id) REFERENCES Boards(id)
-- );

INSERT INTO Panels(id, title, position, board_id) VALUES
    (1, 'Lorem Ipsum', 1, 1),
    (2, 'Lorem Ipsum', 2, 1),
    (3, 'Lorem Ipsum', 3, 1),
    (4, 'Lorem Ipsum', 4, 1),
    (5, 'Lorem Ipsum', 5, 1),
    (6, 'Lorem Ipsum', 6, 1),
    (7, 'Lorem Ipsum 2: Electric Boogaloo', 1, 2),
    (8, 'Lorem Ipsum 2: Electric Boogaloo', 2, 2),
    (9, 'Lorem Ipsum 2: Electric Boogaloo', 3, 2),
    (10, 'Lorem Ipsum 2: Electric Boogaloo', 4, 2),
    (11, 'Lorem Ipsum 2: Electric Boogaloo', 5, 2),
    (12, 'Lorem Ipsum 2: Electric Boogaloo', 6, 2);

-- CREATE TABLE IF NOT EXISTS Stacks (
--     id              SERIAL PRIMARY KEY,
--     title           VARCHAR(255) NOT NULL,
--     position        SMALLINT, -- changed from int to smallint (will not have more than 32767 lists/cards/etc)
--     panel_id        SERIAL, FOREIGN KEY (panel_id) REFERENCES Panels(id)
-- );

INSERT INTO Stacks(id, title, position, panel_id) VALUES
    (1, 'Lorem Ipsum', 1, 1),
    (2, 'Lorem Ipsum', 1, 2),
    (3, 'Lorem Ipsum', 1, 3),
    (4, 'Lorem Ipsum', 1, 4),
    (5, 'Lorem Ipsum', 1, 5),
    (6, 'Lorem Ipsum', 1, 6),
    (7, 'Lorem Ipsum 2: Electric Boogaloo', 1, 7),
    (8, 'Lorem Ipsum 2: Electric Boogaloo', 1, 8),
    (9, 'Lorem Ipsum 2: Electric Boogaloo', 1, 9),
    (10, 'Lorem Ipsum 2: Electric Boogaloo', 1, 10),
    (11, 'Lorem Ipsum 2: Electric Boogaloo', 2, 1),
    (12, 'Lorem Ipsum 3: Return of the Ipsum', 3, 1);

-- CREATE TABLE IF NOT EXISTS Cards (
--     id              SERIAL PRIMARY KEY,
--     title           VARCHAR(255) NOT NULL,
--     description     TEXT,
--     position        SMALLINT, -- changed from int to smallint (will not have more than 32767 lists/cards/etc)
--     stack_id         SERIAL, FOREIGN KEY (stack_id) REFERENCES Stacks(id)
-- );
INSERT INTO Cards(id, title, description, position, stack_id) VALUES
    (1, 'Testing Card', 'A filler card created for demo', 1, 1),
    (2, 'Testing Card', 'A filler card created for demo', 2, 1),
    (3, 'Testing Card', 'A filler card created for demo', 3, 1),
    (4, 'Testing Card', 'A filler card created for demo', 4, 1),
    (5, 'Testing Card', 'A filler card created for demo', 5, 1),
    (6, 'Testing Card', 'A filler card created for demo', 6, 1);

-- CREATE TABLE IF NOT EXISTS Assigned_Cards (
--     id              SERIAL PRIMARY KEY, -- changed from mere join table structure to give an ID for future ease of reference
--     user_id         VARCHAR(64),
--     card_id         SERIAL, FOREIGN KEY (card_id) REFERENCES Cards(id)
-- );

-- CREATE TABLE IF NOT EXISTS Tags (
--     id              SERIAL PRIMARY KEY,
--     name            VARCHAR(255) NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS Card_Tags (
--     id              SERIAL PRIMARY KEY, -- changed from mere join table structure to give an ID for future ease of reference
--     tag_id          SERIAL, FOREIGN KEY (tag_id) REFERENCES Tags(id),
--     card_id         SERIAL, FOREIGN KEY (card_id) REFERENCES Cards(id)
-- );