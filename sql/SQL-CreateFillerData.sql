-- SyncSpace SQL Filler Data Fill

-- Created for https://github.com/Sync-Space-49/syncspace-server/

-- Authors:
--     Nathan Laney
--     Dylan Halstead

SET TIMEZONE = "America/New_York"; -- sets timezone for timestamptz datatypes

INSERT INTO Organizations(id, owner_id, name, description, ai_enabled) VALUES
    ('69639276-a0bb-4a9c-b0b8-9b66b7aada4f', 'auth0|65296d21ab01a819c3034545', 'SyncSpace', 'This is the SyncSpace example organization description. This is a temporary description in filler data. ', true);

INSERT INTO Organizations(id, owner_id, name, description) VALUES
    ('060599b3-2893-4f43-8ce9-9595d6e93f37', 'auth0|65296d21ab01a819c3034545', 'BioQuest', 'Lorem ipsum, dolor sit amet consectetur adipisicing elit. Molestias aut, repellat ipsum facere voluptate dicta obcaecati deserunt nobis suscipit eaque?'),
    ('b6d82b16-2f89-4aa5-b73e-2a826aa77014', 'auth0|65296d21ab01a819c3034545', 'Some Recipe App', 'Lorem ipsum, vulputate eu scelerisque felis imperdiet proin fermentum leo vel orci porta non pulvinar neque laoreet.');

INSERT INTO Boards(id, title, organization_id, owner_id) VALUES
    ('db98f8ee-23b1-48e4-a913-19adc7af59ac', 'SyncSpace Main Board', '69639276-a0bb-4a9c-b0b8-9b66b7aada4f', 'auth0|65296d21ab01a819c3034545', 1),
    ('8c8b39ef-aa76-48c2-a0c6-f53b8debd999', 'A second, cooler board', '69639276-a0bb-4a9c-b0b8-9b66b7aada4f', 'auth0|65296d21ab01a819c3034545'),
    ('42708b4a-1e15-4ae5-b158-861f4a928d9c', 'BioQuest Main Board', '060599b3-2893-4f43-8ce9-9595d6e93f37', 'auth0|65412afb7c403dde6a228283'),
    ('ee8378ef-60f2-4baf-aaff-33a5c395ec5f', 'A second, cooler board', '060599b3-2893-4f43-8ce9-9595d6e93f37', 'auth0|65412afb7c403dde6a228283'),
    ('a07ee223-79ea-4c13-9cda-104baded5ae5', 'Board Title', 'b6d82b16-2f89-4aa5-b73e-2a826aa77014', 'google-oauth2|111709628753664477473'),
    ('7e7f9a4d-0279-4bab-a528-00781fc5b3ca', 'An additional, even COOLER board', '69639276-a0bb-4a9c-b0b8-9b66b7aada4f', 'auth0|65296d21ab01a819c3034545');

INSERT INTO Panels(id, title, position, board_id) VALUES
    ('69639276-a0bb-4a9c-b0b8-9b66b7aaaf31', 'Lorem Ipsum', 1, 'db98f8ee-23b1-48e4-a913-19adc7af59ac'),
    ('69639276-a0bb-4a9c-b0b8-9b66b7aaaf32', 'Lorem Ipsum', 2, 'db98f8ee-23b1-48e4-a913-19adc7af59ac'),
    ('69639276-a0bb-4a9c-b0b8-9b66b7aaaf33', 'Lorem Ipsum', 3, 'db98f8ee-23b1-48e4-a913-19adc7af59ac'),
    ('69639276-a0bb-4a9c-b0b8-9b66b7aaaf34', 'Lorem Ipsum', 4, 'db98f8ee-23b1-48e4-a913-19adc7af59ac'),
    ('69639276-a0bb-4a9c-b0b8-9b66b7aaaf35', 'Lorem Ipsum', 5, 'db98f8ee-23b1-48e4-a913-19adc7af59ac'),
    ('69639276-a0bb-4a9c-b0b8-9b66b7aaaf36', 'Lorem Ipsum', 6, 'db98f8ee-23b1-48e4-a913-19adc7af59ac'),
    ('69639276-a0bb-4a9c-b0b8-9b66b7aaaa31', 'Lorem Ipsum 2: Electric Boogaloo', 1, '8c8b39ef-aa76-48c2-a0c6-f53b8debd999'),
    ('69639276-a0bb-4a9c-b0b8-9b66b7aaab31', 'Lorem Ipsum 2: Electric Boogaloo', 2, '8c8b39ef-aa76-48c2-a0c6-f53b8debd999'),
    ('69639276-a0bb-4a9c-b0b8-9b66b7aaac31', 'Lorem Ipsum 2: Electric Boogaloo', 3, '8c8b39ef-aa76-48c2-a0c6-f53b8debd999'),
    ('69639276-a0bb-4a9c-b0b8-9b66b7aaad31', 'Lorem Ipsum 2: Electric Boogaloo', 4, '8c8b39ef-aa76-48c2-a0c6-f53b8debd999'),
    ('69639276-a0bb-4a9c-b0b8-9b66b7aaaa36', 'Lorem Ipsum 2: Electric Boogaloo', 5, '8c8b39ef-aa76-48c2-a0c6-f53b8debd999'),
    ('69639276-a0bb-4a9c-b0b8-9b66b7aaaf37', 'Lorem Ipsum 2: Electric Boogaloo', 6, '8c8b39ef-aa76-48c2-a0c6-f53b8debd999');

INSERT INTO Stacks(id, title, position, panel_id) VALUES
    ('865be2f9-2025-4a4a-a62f-96509092bea2', 'Lorem Ipsum', 1, '69639276-a0bb-4a9c-b0b8-9b66b7aaaf31'),
    ('772ba005-c5dd-4abd-bcdb-b384a6eb2a8f', 'Lorem Ipsum', 1, '69639276-a0bb-4a9c-b0b8-9b66b7aaaf32'),
    ('4ad3a5cd-0c8e-4689-8102-a13cedb507df', 'Lorem Ipsum', 1, '69639276-a0bb-4a9c-b0b8-9b66b7aaaf33'),
    ('4b7d024b-19de-4405-bb67-950cff4ba30b', 'Lorem Ipsum', 1, '69639276-a0bb-4a9c-b0b8-9b66b7aaaf34'),
    ('3b39ac4c-6d75-46c8-943c-3f1919bc92e0', 'Lorem Ipsum', 1, '69639276-a0bb-4a9c-b0b8-9b66b7aaaf35'),
    ('19aaa28c-da98-46ac-92d8-eac213a262f0', 'Lorem Ipsum', 1, '69639276-a0bb-4a9c-b0b8-9b66b7aaaf36'),
    ('8651d98f-a046-4994-b43c-16baeeee39cf', 'Lorem Ipsum 2: Electric Boogaloo', 1, '69639276-a0bb-4a9c-b0b8-9b66b7aaaf37'),
    ('274e5674-f799-4770-b96f-29915cc8a829', 'Lorem Ipsum 2: Electric Boogaloo', 1, '69639276-a0bb-4a9c-b0b8-9b66b7aaaa36'),
    ('d97574e9-bc25-4178-af45-4487ffae0026', 'Lorem Ipsum 2: Electric Boogaloo', 1, '69639276-a0bb-4a9c-b0b8-9b66b7aaab31'),
    ('53983ff8-cb27-404e-9818-ce364c7493b9', 'Lorem Ipsum 2: Electric Boogaloo', 1, '69639276-a0bb-4a9c-b0b8-9b66b7aaad31'),
    ('fe031a7e-7529-48a2-b674-0447b8e5844f', 'Lorem Ipsum 2: Electric Boogaloo', 2, '69639276-a0bb-4a9c-b0b8-9b66b7aaaf31'),
    ('f4d08871-f86b-4716-bc91-1807ab5efabd', 'Lorem Ipsum 3: Return of the Ipsum', 3, '69639276-a0bb-4a9c-b0b8-9b66b7aaaf31');

INSERT INTO Cards(id, title, description, points, position, stack_id) VALUES
    ('f2785b0b-333a-4cfa-abbc-a668246cebcf', 'Testing Card', 'A filler card created for demo', 'XL', 1, '865be2f9-2025-4a4a-a62f-96509092bea2'),
    ('6f45243d-dd8b-4494-ae6e-96a50cc3dcbf', 'Testing Card', 'A filler card created for demo', 'XS', 2, '865be2f9-2025-4a4a-a62f-96509092bea2'),
    ('4ae5b20d-289f-4a76-9777-2ec727be8379', 'Testing Card', 'A filler card created for demo', 'L', 3, '865be2f9-2025-4a4a-a62f-96509092bea2'),
    ('78c0bc19-cd2c-48df-9a6c-6a6f6f1479c4', 'Testing Card', 'A filler card created for demo', 'M', 4, '865be2f9-2025-4a4a-a62f-96509092bea2'),
    ('61969f1a-0c32-4c6a-9372-1c87c7f9c196', 'Testing Card', 'A filler card created for demo', 'M', 5, '865be2f9-2025-4a4a-a62f-96509092bea2'),
    ('b73e47ba-61ed-4ee6-8806-c95b0d48e7e2', 'Testing Card', 'A filler card created for demo', 'S', 6, '865be2f9-2025-4a4a-a62f-96509092bea2');