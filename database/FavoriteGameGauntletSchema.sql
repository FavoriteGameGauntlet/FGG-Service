--
-- File generated with SQLiteStudio v3.4.17 on Mon Sep 22 00:22:42 2025
--
-- Text encoding used: System
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: EffectHistory
DROP TABLE IF EXISTS EffectHistory;

CREATE TABLE EffectHistory (
    Id             INTEGER PRIMARY KEY AUTOINCREMENT
                           UNIQUE
                           NOT NULL,
    UserId         INTEGER REFERENCES Users (Id) 
                           NOT NULL,
    GameId         INTEGER REFERENCES Games (Id) 
                           NOT NULL,
    ReceivedDate   TEXT    NOT NULL
                           DEFAULT (datetime('now') ),
    RolledDate     TEXT,
    RolledEffectId INTEGER REFERENCES RollEffects (Id) 
);


-- Table: GameHistory
DROP TABLE IF EXISTS GameHistory;

CREATE TABLE GameHistory (
    Id     INTEGER PRIMARY KEY AUTOINCREMENT
                   NOT NULL
                   UNIQUE,
    UserId INTEGER REFERENCES Users (Id),
    GameId INTEGER REFERENCES Games (Id) 
                   NOT NULL,
    State  TEXT    NOT NULL
                   DEFAULT ('rolled'),
    Date   TEXT    NOT NULL
                   DEFAULT (datetime('now') ),
    Result TEXT
);


-- Table: Games
DROP TABLE IF EXISTS Games;

CREATE TABLE Games (
    Id   INTEGER PRIMARY KEY AUTOINCREMENT
                 UNIQUE
                 NOT NULL,
    Name TEXT    NOT NULL,
    Link TEXT
);


-- Table: RollEffects
DROP TABLE IF EXISTS RollEffects;

CREATE TABLE RollEffects (
    Id          INTEGER PRIMARY KEY
                        NOT NULL
                        UNIQUE,
    Name        TEXT    NOT NULL,
    Description TEXT    NOT NULL
);


-- Table: TimerActions
DROP TABLE IF EXISTS TimerActions;

CREATE TABLE TimerActions (
    Id      INTEGER PRIMARY KEY AUTOINCREMENT
                    UNIQUE
                    NOT NULL,
    TimerId INTEGER REFERENCES Timers (Id) 
                    NOT NULL,
    Action  TEXT    NOT NULL,
    Date    TEXT    NOT NULL
                    DEFAULT (datetime('now') ) 
);


-- Table: Timers
DROP TABLE IF EXISTS Timers;

CREATE TABLE Timers (
    Id          INTEGER PRIMARY KEY AUTOINCREMENT
                        UNIQUE
                        NOT NULL,
    UserId      INTEGER REFERENCES Users (Id) 
                        NOT NULL,
    GameId      INTEGER REFERENCES Games (Id),
    State       TEXT    NOT NULL
                        DEFAULT ('created'),
    DurationInS INTEGER NOT NULL
                        DEFAULT (0) 
);


-- Table: UnplayedGames
DROP TABLE IF EXISTS UnplayedGames;

CREATE TABLE UnplayedGames (
    Id     INTEGER PRIMARY KEY AUTOINCREMENT
                   UNIQUE
                   NOT NULL,
    UserId INTEGER REFERENCES Users (Id) 
                   NOT NULL,
    GameId INTEGER REFERENCES Games (Id) 
                   NOT NULL
);


-- Table: Users
DROP TABLE IF EXISTS Users;

CREATE TABLE Users (
    Id   INTEGER PRIMARY KEY AUTOINCREMENT
                 UNIQUE
                 NOT NULL,
    Name TEXT    NOT NULL
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
