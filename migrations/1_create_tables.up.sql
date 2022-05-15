CREATE TABLE IF NOT EXISTS tournaments (
    id INTEGER PRIMARY KEY autoincrement,
    type TEXT,
    name TEXT,
    date TEXT, -- ISO8601
    contents BLOB,
);

CREATE TABLE IF NOT EXISTS players (
    id INTEGER PRIMARY KEY autoincrement,
    name TEXT,
    team TEXT
);

CREATE TABLE IF NOT EXISTS player_tournaments (
    player_id INTEGER,
    tournament_id INTEGER,
    FOREIGN KEY (player_id) REFERENCES players(id),
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id)
);

CREATE TABLE IF NOT EXISTS player_aliases (
    original_player TEXT
    alias TEXT
);

CREATE INDEX IF NOT EXISTS date_index ON tournaments(date);
CREATE UNIQUE INDEX IF NOT EXISTS name_index ON tournaments(name, type);