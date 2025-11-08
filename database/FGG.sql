--
-- File generated with SQLiteStudio v3.4.17 on Sat Nov 8 22:29:34 2025
--
-- Text encoding used: System
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: EffectHistory
CREATE TABLE IF NOT EXISTS EffectHistory (
    Id             TEXT    PRIMARY KEY
                           UNIQUE
                           NOT NULL,
    UserId         TEXT    REFERENCES Users (Id) 
                           NOT NULL,
    GameId         TEXT    REFERENCES Games (Id),
    CreateDate     TEXT    NOT NULL
                           DEFAULT (datetime('now', 'subsec') ),
    RollDate       TEXT,
    RolledEffectId INTEGER REFERENCES Effects (Id) 
);


-- Table: Effects
CREATE TABLE IF NOT EXISTS Effects (
    Id          TEXT PRIMARY KEY
                     NOT NULL
                     UNIQUE,
    Name        TEXT NOT NULL,
    Description TEXT NOT NULL
);


-- Table: GameHistory
CREATE TABLE IF NOT EXISTS GameHistory (
    Id           TEXT    PRIMARY KEY
                         NOT NULL
                         UNIQUE,
    UserId       TEXT    REFERENCES Users (Id) 
                         NOT NULL,
    GameId       TEXT    REFERENCES Games (Id) 
                         NOT NULL,
    State        TEXT    NOT NULL
                         DEFAULT ('started'),
    CreateDate   TEXT    NOT NULL
                         DEFAULT (datetime('now', 'subsec') ),
    ResultPoints INTEGER,
    FinishDate   TEXT
);


-- Table: Games
CREATE TABLE IF NOT EXISTS Games (
    Id         TEXT PRIMARY KEY
                    UNIQUE
                    NOT NULL,
    Name       TEXT NOT NULL
                    UNIQUE,
    Link       TEXT,
    CreateDate TEXT DEFAULT (datetime('now', 'subsec') ) 
                    NOT NULL
);


-- Table: TimerActions
CREATE TABLE IF NOT EXISTS TimerActions (
    Id               TEXT    PRIMARY KEY
                             UNIQUE
                             NOT NULL,
    TimerId          TEXT    REFERENCES Timers (Id) 
                             NOT NULL,
    Action           TEXT    NOT NULL,
    Date             TEXT    NOT NULL
                             DEFAULT (datetime('now', 'subsec') ),
    RemainingTimeInS INTEGER
);


-- Table: Timers
CREATE TABLE IF NOT EXISTS Timers (
    Id          TEXT    PRIMARY KEY
                        UNIQUE
                        NOT NULL,
    UserId      TEXT    REFERENCES Users (Id) 
                        NOT NULL,
    GameId      TEXT    REFERENCES Games (Id) 
                        NOT NULL,
    State       TEXT    NOT NULL
                        DEFAULT ('created'),
    DurationInS INTEGER NOT NULL
                        DEFAULT (10800),
    CreateDate          NOT NULL
                        DEFAULT (datetime('now', 'subsec') ) 
);


-- Table: UnplayedGames
CREATE TABLE IF NOT EXISTS UnplayedGames (
    Id         TEXT PRIMARY KEY
                    UNIQUE
                    NOT NULL,
    UserId     TEXT REFERENCES Users (Id) 
                    NOT NULL,
    GameId     TEXT REFERENCES Games (Id) 
                    NOT NULL,
    CreateDate      NOT NULL
                    DEFAULT (datetime('now', 'subsec') ) 
);


-- Table: Users
CREATE TABLE IF NOT EXISTS Users (
    Id       TEXT PRIMARY KEY
                  UNIQUE
                  NOT NULL,
    Name     TEXT NOT NULL
                  UNIQUE,
    JoinDate TEXT DEFAULT (datetime('now', 'subsec') ) 
                  NOT NULL
);


-- Trigger: UpdateTimerState
CREATE TRIGGER IF NOT EXISTS UpdateTimerState
                       AFTER INSERT
                          ON TimerActions
                    FOR EACH ROW
BEGIN
    UPDATE Timers
       SET State = CASE new.Action WHEN 'start' THEN 'running' WHEN 'pause' THEN 'paused' WHEN 'stop' THEN 'finished' END
     WHERE Id = new.TimerId;
END;


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
