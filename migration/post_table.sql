CREATE TABLE post (
   id          INTEGER PRIMARY KEY,
   header      TEXT    NOT NULL,
   description TEXT    NOT NULL,
   user_id     INTEGER REFERENCES user (id),
   created_at  INTEGER
);