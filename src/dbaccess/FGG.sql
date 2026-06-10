-- PostgreSQL schema for FGG-Service

CREATE TABLE IF NOT EXISTS Users (
    Id          SERIAL PRIMARY KEY,
    Login       TEXT NOT NULL UNIQUE,
    DisplayName TEXT UNIQUE,
    Email       TEXT NOT NULL UNIQUE,
    Password    TEXT NOT NULL,
    JoinDate    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS UserSessions (
    Id         TEXT PRIMARY KEY,
    UserId     INTEGER NOT NULL REFERENCES Users (Id),
    CreateDate TIMESTAMP NOT NULL DEFAULT NOW(),
    ExpiryDate TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '1 day'
);

CREATE TABLE IF NOT EXISTS UserStats (
    Id               SERIAL PRIMARY KEY,
    UserId           INTEGER NOT NULL UNIQUE REFERENCES Users (Id),
    AvailableRolls   INTEGER NOT NULL DEFAULT 0 CHECK (AvailableRolls >= 0),
    TerritoryHours   INTEGER NOT NULL DEFAULT 0 CHECK (TerritoryHours >= 0),
    ExperiencePoints INTEGER NOT NULL DEFAULT 0 CHECK (ExperiencePoints >= 0),
    TerritoryPoints  INTEGER NOT NULL DEFAULT 0 CHECK (TerritoryPoints >= 0),
    FreePoints       INTEGER NOT NULL DEFAULT 0 CHECK (FreePoints >= 0)
);

CREATE TABLE IF NOT EXISTS Games (
    Id         SERIAL PRIMARY KEY,
    Name       TEXT NOT NULL UNIQUE,
    CreateDate TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS GameHistory (
    Id         SERIAL PRIMARY KEY,
    UserId     INTEGER NOT NULL REFERENCES Users (Id),
    GameId     INTEGER NOT NULL REFERENCES Games (Id),
    State      TEXT NOT NULL DEFAULT 'started',
    ChangeDate TIMESTAMP NOT NULL DEFAULT NOW(),
    FinishDate TIMESTAMP
);

CREATE TABLE IF NOT EXISTS UnplayedGames (
    Id         SERIAL PRIMARY KEY,
    UserId     INTEGER NOT NULL REFERENCES Users (Id),
    GameId     INTEGER NOT NULL REFERENCES Games (Id),
    CreateDate TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS Timers (
    Id               SERIAL PRIMARY KEY,
    UserId           INTEGER NOT NULL REFERENCES Users (Id),
    GameId           INTEGER NOT NULL REFERENCES Games (Id),
    State            TEXT NOT NULL DEFAULT 'created',
    DurationInS      INTEGER NOT NULL,
    RemainingTimeInS INTEGER NOT NULL,
    CreateDate       TIMESTAMP NOT NULL DEFAULT NOW(),
    LastActionDate   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS WheelEffects (
    Id                      SERIAL PRIMARY KEY,
    Name                    TEXT NOT NULL,
    Description             TEXT NOT NULL,
    IsCompleted             INTEGER NOT NULL DEFAULT 0,
    OwnerPointChangeFormula TEXT NOT NULL DEFAULT 'output 0',
    EffectRerollFormula     TEXT NOT NULL DEFAULT 'output 0',
    IsItem                  INTEGER NOT NULL DEFAULT 0,
    IsEffectChoice          INTEGER NOT NULL DEFAULT 0,
    RepeatCount             INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS LastWheelEffects (
    Id            SERIAL PRIMARY KEY,
    UserId        INTEGER NOT NULL REFERENCES Users (Id),
    WheelEffectId INTEGER NOT NULL REFERENCES WheelEffects (Id),
    Position      INTEGER NOT NULL,
    IsApplied     INTEGER NOT NULL DEFAULT 0 CHECK (IsApplied IN (0, 1)),
    RollDate      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS SystemParameters (
    Id    SERIAL PRIMARY KEY,
    Name  TEXT NOT NULL UNIQUE,
    Value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS TerritoryPointHistory (
    Id                SERIAL PRIMARY KEY,
    UserId            INTEGER NOT NULL REFERENCES Users (Id),
    SourceUserId      INTEGER NOT NULL REFERENCES Users (Id),
    ChangeSource      TEXT NOT NULL,
    ChangeValue       INTEGER NOT NULL,
    ActualChangeValue INTEGER NOT NULL,
    FinalValue        INTEGER NOT NULL CHECK (FinalValue >= 0),
    ChangeDate        TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS WheelEffectHistory (
    Id            SERIAL PRIMARY KEY,
    UserId        INTEGER NOT NULL REFERENCES Users (Id),
    WheelEffectId INTEGER NOT NULL REFERENCES WheelEffects (Id),
    RollDate      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS FreePointHistory (
    Id                SERIAL PRIMARY KEY,
    UserId            INTEGER NOT NULL REFERENCES Users (Id),
    SourceUserId      INTEGER NOT NULL REFERENCES Users (Id),
    ChangeSource      TEXT NOT NULL,
    ChangeValue       INTEGER NOT NULL,
    ActualChangeValue INTEGER NOT NULL,
    FinalValue        INTEGER NOT NULL CHECK (FinalValue >= 0),
    WheelEffectId     INTEGER REFERENCES WheelEffectHistory (Id),
    ChangeDate        TIMESTAMP NOT NULL DEFAULT NOW()
);
