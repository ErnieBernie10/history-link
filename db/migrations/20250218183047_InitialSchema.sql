-- migrate:up
-- Create the subjects table
CREATE TABLE subjects (
    id TEXT PRIMARY KEY,
    title VARCHAR NOT NULL,
    summary VARCHAR NOT NULL,
    subject_type VARCHAR(50),
    url VARCHAR NOT NULL,
    weight INTEGER,
    from_date DATE NOT NULL,
    until_date DATE NOT NULL
);
-- Create the subject_relations table
CREATE TABLE subject_relations (
    subject_1 TEXT,
    subject_2 TEXT,
    PRIMARY KEY (subject_1, subject_2),
    FOREIGN KEY (subject_1) REFERENCES subjects(id),
    FOREIGN KEY (subject_2) REFERENCES subjects(id)
);
-- Create the impacts table
CREATE TABLE impacts (
    id TEXT PRIMARY KEY,
    subject_id TEXT NOT NULL,
    reasoning VARCHAR NOT NULL,
    category VARCHAR(50) NOT NULL,
    value INTEGER NOT NULL,
    FOREIGN KEY (subject_id) REFERENCES subjects(id)
);
-- Create the impact_revisions table
CREATE TABLE impact_revisions (
    id TEXT PRIMARY KEY,
    impact_id TEXT NOT NULL,
    FOREIGN KEY (impact_id) REFERENCES impacts(id)
);
-- migrate:down
DROP TABLE IF EXISTS impact_revisions;
DROP TABLE IF EXISTS impacts;
DROP TABLE IF EXISTS subject_relations;
DROP TABLE IF EXISTS subjects;

