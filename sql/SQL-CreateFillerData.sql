-- SyncSpace SQL Filler Data Fill

-- Created for https://github.com/Sync-Space-49/syncspace-server/

-- Authors: 
--     Nathan Laney
--     Dylan Halstead

SET TIMEZONE = "America/New_York"; -- sets timezone for timestamptz datatypes

INSERT INTO Organizations(id, owner_id, name, description) VALUES
    ('69639276-a0bb-4a9c-b0b8-9b66b7aada4f', 'auth0|65296d21ab01a819c3034545', 'SyncSpace', 'This is the SyncSpace example organization description. This is a temporary description in filler data. '),
    ('060599b3-2893-4f43-8ce9-9595d6e93f37', 'auth0|65296d21ab01a819c3034545', 'BioQuest', 'Lorem ipsum, dolor sit amet consectetur adipisicing elit. Molestias aut, repellat ipsum facere voluptate dicta obcaecati deserunt nobis suscipit eaque?'),
    ('b6d82b16-2f89-4aa5-b73e-2a826aa77014', 'auth0|65296d21ab01a819c3034545', 'Some Recipe App', 'Lorem ipsum, vulputate eu scelerisque felis imperdiet proin fermentum leo vel orci porta non pulvinar neque laoreet.');

INSERT INTO Boards(id, title, organization_id) VALUES
    (1, 'SyncSpace Main Board', '69639276-a0bb-4a9c-b0b8-9b66b7aada4f'),
    (2, 'A second, cooler board', '69639276-a0bb-4a9c-b0b8-9b66b7aada4f'),
    (3, 'BioQuest Main Board', '060599b3-2893-4f43-8ce9-9595d6e93f37'),
    (4, 'A second, cooler board', '060599b3-2893-4f43-8ce9-9595d6e93f37'),
    (5, 'Board Title', 'b6d82b16-2f89-4aa5-b73e-2a826aa77014'),
    (6, 'An additional, even COOLER board', '69639276-a0bb-4a9c-b0b8-9b66b7aada4f');

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

INSERT INTO Cards(id, title, description, position, stack_id) VALUES
    (1, 'Testing Card', 'A filler card created for demo', 1, 1),
    (2, 'Testing Card', 'A filler card created for demo', 2, 1),
    (3, 'Testing Card', 'A filler card created for demo', 3, 1),
    (4, 'Testing Card', 'A filler card created for demo', 4, 1),
    (5, 'Testing Card', 'A filler card created for demo', 5, 1),
    (6, 'Testing Card', 'A filler card created for demo', 6, 1);