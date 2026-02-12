--
-- File generated with SQLiteStudio v3.4.17 on Fri Feb 13 00:10:14 2026
--
-- Text encoding used: System
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: AvailableRolls
CREATE TABLE IF NOT EXISTS AvailableRolls (
    Id         INTEGER PRIMARY KEY AUTOINCREMENT
                       UNIQUE
                       NOT NULL,
    UserId     INTEGER REFERENCES Users (Id) 
                       NOT NULL,
    CreateDate TEXT    NOT NULL
                       DEFAULT (datetime('now', 'subsec') ) 
);


-- Table: EffectHistory
CREATE TABLE IF NOT EXISTS EffectHistory (
    Id         INTEGER PRIMARY KEY AUTOINCREMENT
                       NOT NULL
                       UNIQUE,
    UserId     INTEGER REFERENCES Users (Id) 
                       NOT NULL,
    EffectId   INTEGER REFERENCES Effects (Id) 
                       NOT NULL,
    CreateDate TEXT    NOT NULL
                       DEFAULT (datetime('now', 'subsec') ) 
);


-- Table: Effects
CREATE TABLE IF NOT EXISTS Effects (
    Id                      INTEGER PRIMARY KEY AUTOINCREMENT
                                    UNIQUE
                                    NOT NULL,
    Name                    TEXT    NOT NULL,
    Description             TEXT    NOT NULL,
    IsCompleted             INTEGER DEFAULT (0) 
                                    NOT NULL,
    OwnerPointChangeFormula TEXT    NOT NULL
                                    DEFAULT ('output 0'),
    EffectRerollFormula     TEXT    NOT NULL
                                    DEFAULT ('output 0'),
    IsItem                  INTEGER DEFAULT (0) 
                                    NOT NULL,
    IsEffectChoice          INTEGER NOT NULL
                                    DEFAULT (0),
    RepeatCount             INTEGER NOT NULL
                                    DEFAULT (1) 
);


-- Table: GameHistory
CREATE TABLE IF NOT EXISTS GameHistory (
    Id         INTEGER PRIMARY KEY AUTOINCREMENT
                       NOT NULL
                       UNIQUE,
    UserId     INTEGER REFERENCES Users (Id) 
                       NOT NULL,
    GameId     INTEGER NOT NULL
                       REFERENCES Games (Id),
    State      TEXT    NOT NULL
                       DEFAULT ('started'),
    CreateDate TEXT    NOT NULL
                       DEFAULT (datetime('now', 'subsec') ),
    FinishDate TEXT
);


-- Table: Games
CREATE TABLE IF NOT EXISTS Games (
    Id         INTEGER PRIMARY KEY AUTOINCREMENT
                       UNIQUE
                       NOT NULL,
    Name       TEXT    NOT NULL
                       UNIQUE,
    CreateDate TEXT    DEFAULT (datetime('now', 'subsec') ) 
                       NOT NULL
);


-- Table: Timers
CREATE TABLE IF NOT EXISTS Timers (
    Id          INTEGER PRIMARY KEY AUTOINCREMENT
                        UNIQUE
                        NOT NULL,
    UserId      INTEGER REFERENCES Users (Id) 
                        NOT NULL,
    GameId      INTEGER NOT NULL
                        REFERENCES Games (Id),
    State       TEXT    NOT NULL
                        DEFAULT ('created'),
    DurationInS INTEGER NOT NULL,
    CreateDate  TEXT    NOT NULL
                        DEFAULT (datetime('now', 'subsec') ) 
);


-- Table: UnplayedGames
CREATE TABLE IF NOT EXISTS UnplayedGames (
    Id         INTEGER PRIMARY KEY AUTOINCREMENT
                       UNIQUE
                       NOT NULL,
    UserId     INTEGER REFERENCES Users (Id) 
                       NOT NULL,
    GameId     INTEGER NOT NULL
                       REFERENCES Games (Id),
    CreateDate TEXT    NOT NULL
                       DEFAULT (datetime('now', 'subsec') ) 
);


-- Table: Users
CREATE TABLE IF NOT EXISTS Users (
    Id       INTEGER PRIMARY KEY AUTOINCREMENT
                     UNIQUE
                     NOT NULL,
    Name     TEXT    NOT NULL
                     UNIQUE,
    Email    TEXT    UNIQUE
                     NOT NULL,
    Password TEXT    NOT NULL,
    JoinDate TEXT    DEFAULT (datetime('now', 'subsec') ) 
                     NOT NULL
);


-- Table: UserSessions
CREATE TABLE IF NOT EXISTS UserSessions (
    Id         TEXT    PRIMARY KEY
                       UNIQUE
                       NOT NULL,
    UserId     INTEGER REFERENCES Users (Id) 
                       NOT NULL,
    CreateDate TEXT    NOT NULL
                       DEFAULT (datetime('now', 'subsec') ),
    ExpiryDate TEXT    NOT NULL
                       DEFAULT (datetime('now', 'subsec', '+1 day') ) 
);


-- Index: UserGameIndex
CREATE INDEX IF NOT EXISTS UserGameIndex ON Timers (
    UserId,
    GameId
);


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
