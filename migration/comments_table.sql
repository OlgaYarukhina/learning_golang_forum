CREATE TABLE comments (
     id         INTEGER PRIMARY KEY,
     comment    TEXT    NOT NULL,
     post_id    INTEGER REFERENCES post (id),
     user_id    INTEGER REFERENCES user (id),
     created_at INTEGER
);