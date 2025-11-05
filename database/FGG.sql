--
-- File generated with SQLiteStudio v3.4.17 on Thu Nov 6 02:06:26 2025
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
    GameId         TEXT    REFERENCES Games (Id) 
                           NOT NULL,
    ReceivedDate   TEXT    NOT NULL
                           DEFAULT (datetime('now') ),
    RolledDate     TEXT,
    RolledEffectId INTEGER REFERENCES RollEffects (Id) 
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
    RolledDate   TEXT    NOT NULL
                         DEFAULT (datetime('now') ),
    ResultPoints INTEGER,
    FinishDate   TEXT
);


-- Table: Games
DROP TABLE IF EXISTS Games;

CREATE TABLE Games (
    Id   TEXT PRIMARY KEY
              UNIQUE
              NOT NULL,
    Name TEXT NOT NULL
              UNIQUE,
    Link TEXT
);


-- Table: RollEffects
DROP TABLE IF EXISTS RollEffects;

CREATE TABLE RollEffects (
    Id          TEXT PRIMARY KEY
                     NOT NULL
                     UNIQUE,
    Name        TEXT NOT NULL,
    Description TEXT NOT NULL
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
                             DEFAULT (datetime('now') ),
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
                        DEFAULT (10800) 
);


-- Table: UnplayedGames
DROP TABLE IF EXISTS UnplayedGames;

CREATE TABLE UnplayedGames (
    Id     TEXT PRIMARY KEY
                UNIQUE
                NOT NULL,
    UserId TEXT REFERENCES Users (Id) 
                NOT NULL,
    GameId TEXT REFERENCES Games (Id) 
                NOT NULL
);


-- Table: Users
DROP TABLE IF EXISTS Users;

CREATE TABLE Users (
    Id   TEXT PRIMARY KEY
              UNIQUE
              NOT NULL,
    Name TEXT NOT NULL
              UNIQUE
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
