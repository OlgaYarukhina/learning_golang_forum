CREATE TABLE comment_like (
    id         INTEGER PRIMARY KEY,
    comment_id INTEGER REFERENCES comments (id),
    user_id    INTEGER REFERENCES user (id)
);
