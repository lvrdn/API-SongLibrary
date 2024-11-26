CREATE TABLE songs (
    "id" serial PRIMARY KEY,
    "song_name" varchar(100) NOT NULL,
    "group_name" varchar(100) NOT NULL,
    "release_date"  date,
    "text" text,
    "link" varchar(255),
    UNIQUE ("song_name", "group_name")
);
