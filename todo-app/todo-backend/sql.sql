-- Create table
CREATE TABLE tasks (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title      TEXT NOT NULL,
    state      INTEGER NOT NULL CHECK (state IN (0, 1))
);


-- Demo Tasks
INSERT INTO tasks (title, state) VALUES
('Buy groceries', 1),
('Finish project documentation', 0),
('Schedule dentist appointment', 0),
('Pay electricity bill', 1),
('Gym workout - Leg day', 1),
('Read 20 pages of a book', 0),
('Update software dependencies', 0),
('Call Amma', 1);
