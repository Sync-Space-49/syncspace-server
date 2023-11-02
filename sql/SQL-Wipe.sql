-- SyncSpace SQL Setup

-- Created for https://github.com/Sync-Space-49/syncspace-server/

-- Authors:
--     Nathan Laney
--     Dylan Halstead

SET TIMEZONE = "America/New_York"; -- sets timezone for timestamptz datatypes

DROP TABLE IF EXISTS Organizations CASCADE;
DROP TABLE IF EXISTS Boards CASCADE;
DROP TABLE IF EXISTS Panels CASCADE;
DROP TABLE IF EXISTS Stacks CASCADE;
DROP TABLE IF EXISTS Cards CASCADE;
DROP TABLE IF EXISTS Assigned_Cards CASCADE;
DROP TABLE IF EXISTS Tags CASCADE;
DROP TABLE IF EXISTS Card_Tags CASCADE;
