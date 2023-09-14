```
SyncSpace SQL Setup

Created for https://github.com/Sync-Space-49/syncspace-server/

Authors: 
    Nathan Laney
    Dylan Halstead

https://user-images.githubusercontent.com/70990184/267409122-ddc40017-4d0a-4e29-a707-4858b364b630.png
``` 
CREATE TABLE IF NOT EXISTS Users (
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    password        VARCHAR(255) NOT NULL,
    pfp_url         VARCHAR(255)
)

