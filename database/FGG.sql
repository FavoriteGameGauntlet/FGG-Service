--
-- File generated with SQLiteStudio v3.4.17 on Fri Nov 7 02:41:26 2025
--
-- Text encoding used: System
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: EffectHistory
DROP TABLE IF EXISTS EffectHistory;

CREATE TABLE EffectHistory (
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
DROP TABLE IF EXISTS Effects;

CREATE TABLE Effects (
    Id          TEXT PRIMARY KEY
                     NOT NULL
                     UNIQUE,
    Name        TEXT NOT NULL,
    Description TEXT NOT NULL
);


-- Table: GameHistory
DROP TABLE IF EXISTS GameHistory;

CREATE TABLE GameHistory (
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
DROP TABLE IF EXISTS Games;

CREATE TABLE Games (
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
DROP TABLE IF EXISTS TimerActions;

CREATE TABLE TimerActions (
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
DROP TABLE IF EXISTS Timers;

CREATE TABLE Timers (
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
DROP TABLE IF EXISTS UnplayedGames;

CREATE TABLE UnplayedGames (
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
DROP TABLE IF EXISTS Users;

CREATE TABLE Users (
    Id       TEXT PRIMARY KEY
                  UNIQUE
                  NOT NULL,
    Name     TEXT NOT NULL
                  UNIQUE,
    JoinDate TEXT DEFAULT (datetime('now', 'subsec') ) 
                  NOT NULL
);


-- Trigger: UpdateTimerState
DROP TRIGGER IF EXISTS UpdateTimerState;
CREATE TRIGGER UpdateTimerState
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
