--
-- File generated with SQLiteStudio v3.4.17 on Sun Feb 15 22:50:27 2026
--
-- Text encoding used: System
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: FreePointHistory
CREATE TABLE IF NOT EXISTS FreePointHistory (
    Id                INTEGER PRIMARY KEY AUTOINCREMENT
                              UNIQUE
                              NOT NULL,
    UserId            NUMERIC REFERENCES Users (Id) 
                              NOT NULL,
    ChangeSource      TEXT    NOT NULL,
    ChangeValue       INTEGER NOT NULL,
    ActualChangeValue INTEGER NOT NULL,
    FinalValue        INTEGER CHECK (FinalValue >= 0) 
                              NOT NULL,
    ChangeDate        TEXT    NOT NULL
                              DEFAULT (datetime('now', 'subsec') ) 
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
    ChangeDate TEXT    NOT NULL
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


-- Table: LastWheelEffects
CREATE TABLE IF NOT EXISTS LastWheelEffects (
    Id            INTEGER PRIMARY KEY AUTOINCREMENT
                          UNIQUE
                          NOT NULL,
    UserId        INTEGER REFERENCES Users (Id) 
                          NOT NULL,
    WheelEffectId INTEGER REFERENCES WheelEffects (Id) 
                          NOT NULL,
    Position      INTEGER NOT NULL,
    IsApplied     INTEGER CHECK (IsApplied IN (0, 1) ) 
                          NOT NULL
                          DEFAULT (0),
    RollDate      TEXT    NOT NULL
                          DEFAULT (datetime('now', 'subsec') ) 
);


-- Table: TerritoryPointHistory
CREATE TABLE IF NOT EXISTS TerritoryPointHistory (
    Id                INTEGER PRIMARY KEY AUTOINCREMENT
                              UNIQUE
                              NOT NULL,
    UserId            INTEGER REFERENCES Users (Id) 
                              NOT NULL,
    ChangeSource      TEXT    NOT NULL,
    ChangeValue       INTEGER NOT NULL,
    ActualChangeValue INTEGER NOT NULL,
    FinalValue        INTEGER NOT NULL
                              CHECK (FinalValue >= 0),
    ChangeDate        TEXT    NOT NULL
                              DEFAULT (datetime('now', 'subsec') ) 
);


-- Table: Timers
CREATE TABLE IF NOT EXISTS Timers (
    Id               INTEGER PRIMARY KEY AUTOINCREMENT
                             UNIQUE
                             NOT NULL,
    UserId           INTEGER REFERENCES Users (Id) 
                             NOT NULL,
    GameId           INTEGER NOT NULL
                             REFERENCES Games (Id),
    State            TEXT    NOT NULL
                             DEFAULT ('created'),
    DurationInS      INTEGER NOT NULL,
    RemainingTimeInS INTEGER NOT NULL,
    CreateDate       TEXT    NOT NULL
                             DEFAULT (datetime('now', 'subsec') ),
    LastActionDate   TEXT    NOT NULL
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
    Id          INTEGER PRIMARY KEY AUTOINCREMENT
                        UNIQUE
                        NOT NULL,
    Login       TEXT    NOT NULL
                        UNIQUE,
    DisplayName TEXT    UNIQUE,
    Email       TEXT    UNIQUE
                        NOT NULL,
    Password    TEXT    NOT NULL,
    JoinDate    TEXT    DEFAULT (datetime('now', 'subsec') ) 
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


-- Table: UserStats
CREATE TABLE IF NOT EXISTS UserStats (
    Id               INTEGER PRIMARY KEY AUTOINCREMENT
                             UNIQUE
                             NOT NULL,
    UserId           INTEGER REFERENCES Users (Id) 
                             NOT NULL
                             UNIQUE,
    AvailableRolls   INTEGER NOT NULL
                             DEFAULT (0) 
                             CHECK (AvailableRolls >= 0),
    TerritoryHours   INTEGER NOT NULL
                             DEFAULT (0) 
                             CHECK (TerritoryHours >= 0),
    ExperiencePoints INTEGER CHECK (ExperiencePoints >= 0) 
                             NOT NULL
                             DEFAULT (0),
    TerritoryPoints  INTEGER CHECK (TerritoryPoints >= 0) 
                             NOT NULL
                             DEFAULT (0),
    FreePoints       INTEGER CHECK (FreePoints >= 0) 
                             NOT NULL
                             DEFAULT (0) 
);


-- Table: WheelEffectHistory
CREATE TABLE IF NOT EXISTS WheelEffectHistory (
    Id            INTEGER PRIMARY KEY AUTOINCREMENT
                          NOT NULL
                          UNIQUE,
    UserId        INTEGER REFERENCES Users (Id) 
                          NOT NULL,
    WheelEffectId INTEGER NOT NULL
                          REFERENCES WheelEffects (Id),
    RollDate      TEXT    NOT NULL
                          DEFAULT (datetime('now', 'subsec') ) 
);


-- Table: WheelEffects
CREATE TABLE IF NOT EXISTS WheelEffects (
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


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
