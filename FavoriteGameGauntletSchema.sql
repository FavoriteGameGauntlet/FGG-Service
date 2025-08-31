--
-- File generated with SQLiteStudio v3.4.17 on Sun Aug 31 23:21:26 2025
--
-- Text encoding used: System
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: RollEffects
CREATE TABLE RollEffects (
    Id          INTEGER PRIMARY KEY
                        NOT NULL
                        UNIQUE,
    Description TEXT    NOT NULL
);


-- Table: RollHistory
CREATE TABLE RollHistory (
    Id             INTEGER PRIMARY KEY AUTOINCREMENT
                           UNIQUE
                           NOT NULL,
    UserId         INTEGER REFERENCES Users (Id) 
                           NOT NULL,
    RecievedDate   TEXT    NOT NULL
                           DEFAULT (datetime('now') ),
    RolledDate     TEXT,
    RolledEffectId INTEGER REFERENCES RollEffects (Id) 
);


-- Table: TimerActions
CREATE TABLE TimerActions (
    Id      INTEGER PRIMARY KEY AUTOINCREMENT
                    UNIQUE
                    NOT NULL,
    TimerId INTEGER REFERENCES Timers (Id) 
                    NOT NULL,
    Action  INTEGER NOT NULL,
    Date    TEXT    NOT NULL
                    DEFAULT (datetime('now') ) 
);


-- Table: Timers
CREATE TABLE Timers (
    Id          INTEGER PRIMARY KEY AUTOINCREMENT
                        UNIQUE
                        NOT NULL,
    UserId      INTEGER REFERENCES Users (Id) 
                        NOT NULL,
    State       INTEGER NOT NULL
                        DEFAULT (0),
    DurationInS INTEGER NOT NULL
                        DEFAULT (0) 
);


-- Table: Users
CREATE TABLE Users (
    Id   INTEGER PRIMARY KEY AUTOINCREMENT
                 UNIQUE
                 NOT NULL,
    Name TEXT    NOT NULL
                 UNIQUE
);


-- Trigger: UpdateTimerState
CREATE TRIGGER UpdateTimerState
         AFTER INSERT
            ON TimerActions
      FOR EACH ROW
BEGIN
    UPDATE Timers
       SET State = CASE new.Action WHEN 0 THEN 1 WHEN 1 THEN 2 WHEN 2 THEN 3 END
     WHERE Id = new.TimerId;
END;


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
