-- Create table
CREATE TABLE tasks (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title      TEXT NOT NULL,
    state      INTEGER NOT NULL CHECK (state IN (0, 1))
);
